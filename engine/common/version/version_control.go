package version

import (
	"errors"
	"fmt"
	"math"
	"sync"

	"github.com/coreos/go-semver/semver"
	"github.com/rs/zerolog"

	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/module/component"
	"github.com/onflow/flow-go/module/irrecoverable"
	"github.com/onflow/flow-go/state/protocol"
	psEvents "github.com/onflow/flow-go/state/protocol/events"
	"github.com/onflow/flow-go/storage"
)

// ErrUnhandledHeight
var ErrUnhandledHeight = errors.New("unhandled height")

// ErrUnhandledVersionBeacon
var ErrUnhandledVersionBeacon = errors.New("unhandled version beacon")

// VersionControlConsumer defines a function type that consumes version control updates.
// It is called with the block height and the corresponding sem-version.
type VersionControlConsumer func(height uint64, semver string)

// NoHeight represents the maximum possible height for blocks.
var NoHeight = uint64(math.MaxUint64)

// VersionControl manages the version control system for the node.
// It consumes BlockFinalized events and updates the node's version control based on the latest version beacon.
//
// VersionControl implements the protocol.Consumer and component.Component interfaces.
type VersionControl struct {
	// Noop implements the protocol.Consumer interface with no operations.
	psEvents.Noop
	sync.RWMutex
	component.Component

	log zerolog.Logger
	// Storage
	versionBeacons storage.VersionBeacons

	// nodeVersion stores the node's current version.
	// It could be nil if the node version is not available.
	nodeVersion *semver.Version
	// versionBeacon stores the last handled version beacon.
	versionBeacon *flow.SealedVersionBeacon

	// consumers stores the list of consumers for version updates.
	consumers []VersionControlConsumer

	// blockFinalizedChan is a channel for receiving block finalized events.
	blockFinalizedChan chan uint64

	// lastHandledBlockHeight the last handled block height
	lastHandledBlockHeight uint64

	// startHeight and endHeight define the height boundaries for version compatibility.
	startHeight uint64
	endHeight   uint64
}

var _ protocol.Consumer = (*VersionControl)(nil)
var _ component.Component = (*VersionControl)(nil)

// NewVersionControl creates a new VersionControl instance.
//
// We currently have no strong guarantee that the node version is a valid semver.
// See build.SemverV2 for more details. That is why nil is a valid input for node version.
// Without a node version, the stop control can still be used for manual stopping.
func NewVersionControl(
	log zerolog.Logger,
	versionBeacons storage.VersionBeacons,
	nodeVersion *semver.Version,
	latestFinalizedBlock *flow.Header,
) *VersionControl {
	// blockFinalizedChan is buffered to handle block finalized events.
	blockFinalizedChan := make(chan uint64, 1000)

	vc := &VersionControl{
		log: log.With().
			Str("component", "version_control").
			Logger(),

		blockFinalizedChan: blockFinalizedChan,
		nodeVersion:        nodeVersion,
		versionBeacons:     versionBeacons,
		startHeight:        NoHeight,
		endHeight:          NoHeight,
	}

	if vc.nodeVersion != nil {
		log = log.With().
			Stringer("node_version", vc.nodeVersion).
			Logger()
	}

	log.Info().Msg("system initialized")

	// Setup component manager for handling worker functions.
	cm := component.NewComponentManagerBuilder()
	cm.AddWorker(vc.processEvents)
	cm.AddWorker(vc.checkInitialVersionBeacon(latestFinalizedBlock.Height))

	vc.Component = cm.Build()

	return vc
}

//TODO: This should go down through all version beacons events

// checkInitialVersionBeacon checks the initial version beacon at the latest finalized block.
// It ensures the component is not ready until the initial version beacon is checked.
func (v *VersionControl) checkInitialVersionBeacon(
	latestFinalizedBlockHeight uint64,
) func(ctx irrecoverable.SignalerContext, ready component.ReadyFunc) {
	return func(ctx irrecoverable.SignalerContext, ready component.ReadyFunc) {
		// component is not ready until we checked the initial version beacon
		defer ready()

		vb, err := v.versionBeacons.Highest(latestFinalizedBlockHeight)
		if err != nil {
			ctx.Throw(
				fmt.Errorf(
					"failed to get highest version beacon for version control: %w",
					err))
			return
		}

		if vb == nil {
			// no version beacon found
			// this is unexpected as there should always be at least the
			// starting version beacon, but not fatal.
			// It can happen if the node starts before bootstrap is finished.
			// TODO: remove when we can guarantee that there will always be a version beacon
			v.log.Info().
				Uint64("height", latestFinalizedBlockHeight).
				Msg("No version beacon found for version control")
			return
		}

		v.versionBeacon = vb
		v.lastHandledBlockHeight = latestFinalizedBlockHeight

		lastCompatibleHeight := v.lastHandledBlockHeight

		// version boundaries are sorted by version
		for _, boundary := range vb.VersionBoundaries {
			ver, err := boundary.Semver()
			if err != nil || ver == nil {
				// this should never happen as we already validated the version beacon
				// when indexing it
				ctx.Throw(
					fmt.Errorf(
						"failed to parse semver during version control setup: %w",
						err))
				return
			}

			if ver.LessThan(*v.nodeVersion) {
				v.log.Info().
					Uint64("startHeight", lastCompatibleHeight).
					Msg("Found start block height for current node version")
				v.startHeight = lastCompatibleHeight
				return
			}

			lastCompatibleHeight = boundary.BlockHeight
		}
	}
}

// BlockFinalized is called when a block is finalized.
// It implements the protocol.Consumer interface.
func (v *VersionControl) BlockFinalized(h *flow.Header) {
	v.blockFinalizedChan <- h.Height
}

// BlockFinalizedForTesting is used for testing purposes only.
// It simulates a block finalized event for testing the version control logic.
func (v *VersionControl) BlockFinalizedForTesting(h *flow.Header) {
	v.blockFinalized(irrecoverable.MockSignalerContext{}, h.Height)
}

// CompatibleAtBlock checks if the node's version is compatible at a given block height.
// It returns true if the node's version is compatible within the specified height range.
// Returns expected errors:
// - ErrUnhandledVersionBeacon if no version beacon was handled yet
// - ErrUnhandledHeight if incoming block height is higher that last handled block height
func (v *VersionControl) CompatibleAtBlock(height uint64) (bool, error) {
	v.RLock()
	defer v.RUnlock()

	// Check if the version beacon is nil. If so, return an error indicating that the version beacon is unhandled.
	if v.versionBeacon == nil {
		return false, fmt.Errorf("could not check compatibility for height %d: %w", height, ErrUnhandledVersionBeacon)
	}

	// Check if the height is greater than the last handled block height. If so, return an error indicating that the height is unhandled.
	if height > v.lastHandledBlockHeight {
		return false, fmt.Errorf("could not check compatibility for height %d: last handled height is %d: %w", height, v.lastHandledBlockHeight, ErrUnhandledHeight)
	}

	// Check if the start height is set and the height is less than the start height. If so, return false indicating that the height is not compatible.
	if v.startHeight != NoHeight && height < v.startHeight {
		return false, nil
	}

	// Check if the end height is set and the height is greater than the end height. If so, return false indicating that the height is not compatible.
	if v.endHeight != NoHeight && height > v.endHeight {
		return false, nil
	}

	// If none of the above conditions are met, the height is compatible.
	return true, nil
}

// AddVersionUpdatesConsumer adds a consumer for version update events.
func (v *VersionControl) AddVersionUpdatesConsumer(consumer VersionControlConsumer) {
	v.Lock()
	defer v.Unlock()

	v.consumers = append(v.consumers, consumer)
}

// processEvents is a worker that processes block finalized events.
func (v *VersionControl) processEvents(
	ctx irrecoverable.SignalerContext,
	ready component.ReadyFunc,
) {
	ready()

	// TODO: Height tracking mechanism

	for {
		select {
		case <-ctx.Done():
			return
		case h := <-v.blockFinalizedChan:
			v.blockFinalized(ctx, h)
		}
	}
}

// blockFinalized processes a block finalized event and updates the version control state.
func (v *VersionControl) blockFinalized(
	ctx irrecoverable.SignalerContext,
	height uint64,
) {
	v.Lock()
	defer v.Unlock()

	if v.versionBeacon != nil && v.versionBeacon.SealHeight >= height {
		// already processed this or a higher version beacon
		return
	}

	vb, err := v.versionBeacons.Highest(height)
	if err != nil {
		v.log.Err(err).
			Uint64("height", height).
			Msg("Failed to get highest version beacon for version control")

		ctx.Throw(
			fmt.Errorf(
				"failed to get highest version beacon for version control: %w",
				err))
		return
	}

	if vb == nil {
		// no version beacon found
		// this is unexpected as there should always be at least the
		// starting version beacon, but not fatal.
		// It can happen if the node starts before bootstrap is finished.
		// TODO: remove when we can guarantee that there will always be a version beacon
		v.log.Info().
			Uint64("height", height).
			Msg("No version beacon found for version control")
		return
	}

	if v.versionBeacon != nil && v.versionBeacon.SealHeight >= vb.SealHeight {
		// we already processed this or a higher version beacon
		return
	}

	v.versionBeacon = vb
	v.lastHandledBlockHeight = height

	// version boundaries are sorted by version
	for _, boundary := range v.versionBeacon.VersionBoundaries {
		ver, err := boundary.Semver()
		if err != nil || ver == nil {
			// this should never happen as we already validated the version beacon
			// when indexing it
			ctx.Throw(
				fmt.Errorf(
					"failed to parse semver: %w",
					err))
			return
		}

		if ver.Compare(*v.nodeVersion) > 0 {
			v.endHeight = boundary.BlockHeight - 1

			for _, consumer := range v.consumers {
				consumer(boundary.BlockHeight, ver.String())
			}

			break
		}
	}
}
