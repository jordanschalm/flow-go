package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/module"
)

type ComplianceCollector struct {
	finalizedHeight            prometheus.Gauge
	sealedHeight               prometheus.Gauge
	finalizedBlocks            prometheus.Counter
	sealedBlocks               prometheus.Counter
	finalizedPayload           *prometheus.CounterVec
	sealedPayload              *prometheus.CounterVec
	lastBlockFinalizedAt       time.Time
	finalizedBlocksPerSecond   prometheus.Summary
	lastEpochTransitionHeight  prometheus.Gauge
	currentEpochCounter        prometheus.Gauge
	currentEpochPhase          prometheus.Gauge
	currentEpochFinalView      prometheus.Gauge
	currentDKGPhase1FinalView  prometheus.Gauge
	currentDKGPhase2FinalView  prometheus.Gauge
	currentDKGPhase3FinalView  prometheus.Gauge
	epochFallbackModeTriggered prometheus.Gauge
	protocolStateVersion       prometheus.Gauge
}

var _ module.ComplianceMetrics = (*ComplianceCollector)(nil)

func NewComplianceCollector() *ComplianceCollector {

	cc := &ComplianceCollector{

		currentEpochCounter: promauto.NewGauge(prometheus.GaugeOpts{
			Name:      "current_epoch_counter",
			Namespace: namespaceConsensus,
			Subsystem: subsystemCompliance,
			Help:      "the current epoch's counter",
		}),

		currentEpochPhase: promauto.NewGauge(prometheus.GaugeOpts{
			Name:      "current_epoch_phase",
			Namespace: namespaceConsensus,
			Subsystem: subsystemCompliance,
			Help:      "the current epoch's phase",
		}),

		lastEpochTransitionHeight: promauto.NewGauge(prometheus.GaugeOpts{
			Name:      "last_epoch_transition_height",
			Namespace: namespaceConsensus,
			Subsystem: subsystemCompliance,
			Help:      "the height of the most recent finalized epoch transition; in other words the height of the first block of the current epoch",
		}),

		currentEpochFinalView: promauto.NewGauge(prometheus.GaugeOpts{
			Name:      "current_epoch_final_view",
			Namespace: namespaceConsensus,
			Subsystem: subsystemCompliance,
			Help:      "the final view of the current epoch",
		}),

		currentDKGPhase1FinalView: promauto.NewGauge(prometheus.GaugeOpts{
			Name:      "current_dkg_phase1_final_view",
			Namespace: namespaceConsensus,
			Subsystem: subsystemCompliance,
			Help:      "the final view of phase 1 of the current epochs DKG",
		}),
		currentDKGPhase2FinalView: promauto.NewGauge(prometheus.GaugeOpts{
			Name:      "current_dkg_phase2_final_view",
			Namespace: namespaceConsensus,
			Subsystem: subsystemCompliance,
			Help:      "the final view of phase 2 of current epochs DKG",
		}),

		currentDKGPhase3FinalView: promauto.NewGauge(prometheus.GaugeOpts{
			Name:      "current_dkg_phase3_final_view",
			Namespace: namespaceConsensus,
			Subsystem: subsystemCompliance,
			Help:      "the final view of phase 3 of the current epochs DKG (a successful DKG will end shortly after this view)",
		}),

		finalizedHeight: promauto.NewGauge(prometheus.GaugeOpts{
			Name:      "finalized_height",
			Namespace: namespaceConsensus,
			Subsystem: subsystemCompliance,
			Help:      "the last finalized height",
		}),

		sealedHeight: promauto.NewGauge(prometheus.GaugeOpts{
			Name:      "sealed_height",
			Namespace: namespaceConsensus,
			Subsystem: subsystemCompliance,
			Help:      "the last sealed height",
		}),

		finalizedBlocks: promauto.NewCounter(prometheus.CounterOpts{
			Name:      "finalized_blocks_total",
			Namespace: namespaceConsensus,
			Subsystem: subsystemCompliance,
			Help:      "the number of finalized blocks",
		}),

		sealedBlocks: promauto.NewCounter(prometheus.CounterOpts{
			Name:      "sealed_blocks_total",
			Namespace: namespaceConsensus,
			Subsystem: subsystemCompliance,
			Help:      "the number of sealed blocks",
		}),

		finalizedPayload: promauto.NewCounterVec(prometheus.CounterOpts{
			Name:      "finalized_payload_total",
			Namespace: namespaceConsensus,
			Subsystem: subsystemCompliance,
			Help:      "the number of resources in finalized blocks",
		}, []string{LabelResource}),

		sealedPayload: promauto.NewCounterVec(prometheus.CounterOpts{
			Name:      "sealed_payload_total",
			Namespace: namespaceConsensus,
			Subsystem: subsystemCompliance,
			Help:      "the number of resources in sealed blocks",
		}, []string{LabelResource}),

		finalizedBlocksPerSecond: promauto.NewSummary(prometheus.SummaryOpts{
			Name:      "finalized_blocks_per_second",
			Namespace: namespaceConsensus,
			Subsystem: subsystemCompliance,
			Help:      "the number of finalized blocks per second/the finalized block rate",
			Objectives: map[float64]float64{
				0.01: 0.001,
				0.1:  0.01,
				0.5:  0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
			MaxAge:     10 * time.Minute,
			AgeBuckets: 5,
			BufCap:     500,
		}),

		epochFallbackModeTriggered: promauto.NewGauge(prometheus.GaugeOpts{
			Name:      "epoch_fallback_triggered",
			Namespace: namespaceConsensus,
			Subsystem: subsystemCompliance,
			Help:      "indicates whether epoch fallback mode is triggered; if >0, the fallback is triggered",
		}),

		protocolStateVersion: promauto.NewGauge(prometheus.GaugeOpts{
			Name:      "protocol_state_version",
			Namespace: namespaceConsensus,
			Subsystem: subsystemCompliance,
			Help:      "reports the protocol state version of the latest finalized block",
		}),
	}

	return cc
}

// FinalizedHeight sets the finalized height.
func (cc *ComplianceCollector) FinalizedHeight(height uint64) {
	cc.finalizedHeight.Set(float64(height))
}

// BlockFinalized reports metrics about finalized blocks.
func (cc *ComplianceCollector) BlockFinalized(block *flow.Block) {
	now := time.Now()
	if !cc.lastBlockFinalizedAt.IsZero() {
		cc.finalizedBlocksPerSecond.Observe(1.0 / now.Sub(cc.lastBlockFinalizedAt).Seconds())
	}
	cc.lastBlockFinalizedAt = now

	cc.finalizedBlocks.Inc()
	cc.finalizedPayload.With(prometheus.Labels{LabelResource: ResourceGuarantee}).Add(float64(len(block.Payload.Guarantees)))
	cc.finalizedPayload.With(prometheus.Labels{LabelResource: ResourceSeal}).Add(float64(len(block.Payload.Seals)))
}

// SealedHeight sets the finalized height.
func (cc *ComplianceCollector) SealedHeight(height uint64) {
	cc.sealedHeight.Set(float64(height))
}

// BlockSealed reports metrics about sealed blocks.
func (cc *ComplianceCollector) BlockSealed(block *flow.Block) {
	cc.sealedBlocks.Inc()
	cc.sealedPayload.With(prometheus.Labels{LabelResource: ResourceGuarantee}).Add(float64(len(block.Payload.Guarantees)))
	cc.sealedPayload.With(prometheus.Labels{LabelResource: ResourceSeal}).Add(float64(len(block.Payload.Seals)))
}

func (cc *ComplianceCollector) EpochTransitionHeight(height uint64) {
	// An epoch transition comprises a block in epoch N followed by a block in epoch N+1.
	// height here refers to the height of the first block in epoch N+1.
	cc.lastEpochTransitionHeight.Set(float64(height))
}

func (cc *ComplianceCollector) CurrentEpochCounter(counter uint64) {
	cc.currentEpochCounter.Set(float64(counter))
}

func (cc *ComplianceCollector) CurrentEpochPhase(phase flow.EpochPhase) {
	cc.currentEpochPhase.Set(float64(phase))
}

func (cc *ComplianceCollector) CurrentEpochFinalView(view uint64) {
	cc.currentEpochFinalView.Set(float64(view))
}

func (cc *ComplianceCollector) CurrentDKGPhaseViews(phase1FinalView, phase2FinalView, phase3FinalView uint64) {
	cc.currentDKGPhase1FinalView.Set(float64(phase1FinalView))
	cc.currentDKGPhase2FinalView.Set(float64(phase2FinalView))
	cc.currentDKGPhase3FinalView.Set(float64(phase3FinalView))
}

func (cc *ComplianceCollector) EpochFallbackModeTriggered() {
	cc.epochFallbackModeTriggered.Set(float64(1))
}

func (cc *ComplianceCollector) EpochFallbackModeExited() {
	cc.epochFallbackModeTriggered.Set(float64(0))
}

func (cc *ComplianceCollector) ProtocolStateVersion(version uint64) {
	cc.protocolStateVersion.Set(float64(version))
}
