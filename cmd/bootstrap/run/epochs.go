package run

import (
	"encoding/hex"
	"fmt"
	"github.com/onflow/crypto"

	"github.com/rs/zerolog"

	"github.com/onflow/cadence"

	"github.com/onflow/flow-go/cmd/util/cmd/common"
	"github.com/onflow/flow-go/fvm/systemcontracts"
	"github.com/onflow/flow-go/model/bootstrap"
	model "github.com/onflow/flow-go/model/bootstrap"
	"github.com/onflow/flow-go/model/cluster"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/model/flow/filter"
	"github.com/onflow/flow-go/state/protocol/inmem"
)

// GenerateRecoverEpochTxArgs generates the required transaction arguments for the `recoverEpoch` transaction.
// No errors are expected during normal operation.
func GenerateRecoverEpochTxArgs(log zerolog.Logger,
	internalNodePrivInfoDir string,
	nodeConfigJson string,
	collectionClusters int,
	recoveryEpochCounter uint64,
	rootChainID flow.ChainID,
	numViewsInStakingAuction uint64,
	numViewsInEpoch uint64,
	targetDuration uint64,
	unsafeAllowOverWrite bool,
	snapshot *inmem.Snapshot,
) ([]cadence.Value, error) {
	epoch := snapshot.Epochs().Current()

	currentEpochIdentities, err := snapshot.Identities(filter.IsValidProtocolParticipant)
	if err != nil {
		return nil, fmt.Errorf("failed to get  valid protocol participants from snapshot: %w", err)
	}
	// We need canonical ordering here; sanity check to enforce this:
	if !currentEpochIdentities.Sorted(flow.Canonical[flow.Identity]) {
		return nil, fmt.Errorf("identies from snapshot not in canonical order")
	}

	// separate collector nodes by internal and partner nodes
	collectors := currentEpochIdentities.Filter(filter.HasRole[flow.Identity](flow.RoleCollection))
	internalCollectors := make(flow.IdentityList, 0)
	partnerCollectors := make(flow.IdentityList, 0)

	log.Info().Msg("collecting internal node network and staking keys")
	internalNodes, err := common.ReadFullInternalNodeInfos(log, internalNodePrivInfoDir, nodeConfigJson)
	if err != nil {
		return nil, fmt.Errorf("failed to read full internal node infos: %w", err)
	}

	internalNodesMap := make(map[flow.Identifier]struct{})
	for _, node := range internalNodes {
		internalNodesMap[node.NodeID] = struct{}{}
	}
	log.Info().Msg("")

	for _, collector := range collectors {
		if _, ok := internalNodesMap[collector.NodeID]; ok {
			internalCollectors = append(internalCollectors, collector)
		} else {
			partnerCollectors = append(partnerCollectors, collector)
		}
	}

	log.Info().Msg("computing collection node clusters")

	assignments, clusters, err := common.ConstructClusterAssignment(log, partnerCollectors, internalCollectors, collectionClusters)
	if err != nil {
		return nil, fmt.Errorf("unable to generate cluster assignment: %w", err)
	}
	log.Info().Msg("")

	log.Info().Msg("constructing root blocks for collection node clusters")
	clusterBlocks := GenerateRootClusterBlocks(recoveryEpochCounter, clusters)
	log.Info().Msg("")

	log.Info().Msg("constructing root QCs for collection node clusters")
	clusterQCs := ConstructRootQCsForClusters(log, clusters, internalNodes, clusterBlocks)
	log.Info().Msg("")

	epochProtocolState, err := snapshot.EpochProtocolState()
	if err != nil {
		return nil, fmt.Errorf("failed to get epoch protocol state from snapshot: %w", err)
	}
	currentEpochCommit := epochProtocolState.EpochCommit()

	// NOTE: The RecoveryEpoch will re-use the last successful DKG output. This means that the random beacon committee can be
	// different from the consensus committee. This could happen if the node was ejected from the consensus committee, but it still has to be
	// included in the DKG committee since the threshold signature scheme operates on pre-defined number of participants and cannot be changed.
	dkgGroupKeyCdc, cdcErr := cadence.NewString(hex.EncodeToString(currentEpochCommit.DKGGroupKey.Encode()))
	if cdcErr != nil {
		return nil, fmt.Errorf("failed to get dkg group key cadence string: %w", cdcErr)
	}

	// copy DKG index map from the current epoch
	dkgIndexMapPairs := make([]cadence.KeyValuePair, 0)
	for nodeID, index := range currentEpochCommit.DKGIndexMap {
		dkgIndexMapPairs = append(dkgIndexMapPairs, cadence.KeyValuePair{
			Key:   cadence.String(nodeID.String()),
			Value: cadence.NewInt(index),
		})
	}
	// copy DKG public keys from the current epoch
	dkgPubKeys := make([]cadence.Value, 0)
	for _, dkgPubKey := range currentEpochCommit.DKGParticipantKeys {
		dkgPubKeyCdc, cdcErr := cadence.NewString(hex.EncodeToString(dkgPubKey.Encode()))
		if cdcErr != nil {
			return nil, fmt.Errorf("failed to get dkg pub key cadence string for node: %w", cdcErr)
		}
		dkgPubKeys = append(dkgPubKeys, dkgPubKeyCdc)
	}
	// fill node IDs
	nodeIds := make([]cadence.Value, 0)
	for _, id := range currentEpochIdentities {
		nodeIdCdc, err := cadence.NewString(id.GetNodeID().String())
		if err != nil {
			return nil, fmt.Errorf("failed to convert node ID to cadence string %s: %w", id.GetNodeID(), err)
		}
		nodeIds = append(nodeIds, nodeIdCdc)
	}

	clusterQCAddress := systemcontracts.SystemContractsForChain(rootChainID).ClusterQC.Address.String()
	qcVoteData, err := common.ConvertClusterQcsCdc(clusterQCs, clusters, clusterQCAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to convert cluster qcs to cadence type")
	}
	currEpochFinalView, err := epoch.FinalView()
	if err != nil {
		return nil, fmt.Errorf("failed to get final view of current epoch")
	}
	currEpochTargetEndTime, err := epoch.TargetEndTime()
	if err != nil {
		return nil, fmt.Errorf("failed to get target end time of current epoch")
	}

	args := []cadence.Value{
		// recovery epoch counter
		cadence.NewUInt64(recoveryEpochCounter),
		// epoch start view
		cadence.NewUInt64(currEpochFinalView + 1),
		// staking phase end view
		cadence.NewUInt64(currEpochFinalView + numViewsInStakingAuction),
		// epoch end view
		cadence.NewUInt64(currEpochFinalView + numViewsInEpoch),
		// target duration
		cadence.NewUInt64(targetDuration),
		// target end time
		cadence.NewUInt64(currEpochTargetEndTime),
		// clusters,
		common.ConvertClusterAssignmentsCdc(assignments),
		// qcVoteData
		cadence.NewArray(qcVoteData),
		// dkg pub keys
		cadence.NewArray(dkgPubKeys),
		// dkg group key,
		dkgGroupKeyCdc,
		// dkg index map
		cadence.NewDictionary(dkgIndexMapPairs),
		// node ids
		cadence.NewArray(nodeIds),
		// recover the network by initializing a new recover epoch which will increment the smart contract epoch counter
		// or overwrite the epoch metadata for the current epoch
		cadence.NewBool(unsafeAllowOverWrite),
	}

	return args, nil
}

func GenerateRecoverTxArgsWithDKG(log zerolog.Logger,
	internalNodes []bootstrap.NodeInfo,
	collectionClusters int,
	recoveryEpochCounter uint64,
	rootChainID flow.ChainID,
	numViewsInStakingAuction uint64,
	numViewsInEpoch uint64,
	targetDuration uint64,
	unsafeAllowOverWrite bool,
	dkgIndexMap flow.DKGIndexMap,
	dkgParticipantKeys []crypto.PublicKey,
	dkgGroupKey crypto.PublicKey,
	snapshot *inmem.Snapshot,
) ([]cadence.Value, error) {
	epoch := snapshot.Epochs().Current()

	currentEpochIdentities, err := snapshot.Identities(filter.IsValidProtocolParticipant)
	if err != nil {
		return nil, fmt.Errorf("failed to get  valid protocol participants from snapshot: %w", err)
	}
	// We need canonical ordering here; sanity check to enforce this:
	if !currentEpochIdentities.Sorted(flow.Canonical[flow.Identity]) {
		return nil, fmt.Errorf("identies from snapshot not in canonical order")
	}

	// separate collector nodes by internal and partner nodes
	collectors := currentEpochIdentities.Filter(filter.HasRole[flow.Identity](flow.RoleCollection))
	internalCollectors := make(flow.IdentityList, 0)
	partnerCollectors := make(flow.IdentityList, 0)

	internalNodesMap := make(map[flow.Identifier]struct{})
	for _, node := range internalNodes {
		internalNodesMap[node.NodeID] = struct{}{}
	}

	for _, collector := range collectors {
		if _, ok := internalNodesMap[collector.NodeID]; ok {
			internalCollectors = append(internalCollectors, collector)
		} else {
			partnerCollectors = append(partnerCollectors, collector)
		}
	}

	assignments, clusters, err := common.ConstructClusterAssignment(log, partnerCollectors, internalCollectors, collectionClusters)
	if err != nil {
		return nil, fmt.Errorf("unable to generate cluster assignment: %w", err)
	}

	clusterBlocks := GenerateRootClusterBlocks(recoveryEpochCounter, clusters)
	clusterQCs := ConstructRootQCsForClusters(log, clusters, internalNodes, clusterBlocks)

	// NOTE: The RecoveryEpoch will re-use the last successful DKG output. This means that the random beacon committee can be
	// different from the consensus committee. This could happen if the node was ejected from the consensus committee, but it still has to be
	// included in the DKG committee since the threshold signature scheme operates on pre-defined number of participants and cannot be changed.
	dkgGroupKeyCdc, cdcErr := cadence.NewString(hex.EncodeToString(dkgGroupKey.Encode()))
	if cdcErr != nil {
		return nil, fmt.Errorf("failed to get dkg group key cadence string: %w", cdcErr)
	}

	// copy DKG index map from the current epoch
	dkgIndexMapPairs := make([]cadence.KeyValuePair, 0)
	for nodeID, index := range dkgIndexMap {
		dkgIndexMapPairs = append(dkgIndexMapPairs, cadence.KeyValuePair{
			Key:   cadence.String(nodeID.String()),
			Value: cadence.NewInt(index),
		})
	}
	// copy DKG public keys from the current epoch
	dkgPubKeys := make([]cadence.Value, 0)
	for _, dkgPubKey := range dkgParticipantKeys {
		dkgPubKeyCdc, cdcErr := cadence.NewString(hex.EncodeToString(dkgPubKey.Encode()))
		if cdcErr != nil {
			return nil, fmt.Errorf("failed to get dkg pub key cadence string for node: %w", cdcErr)
		}
		dkgPubKeys = append(dkgPubKeys, dkgPubKeyCdc)
	}
	// fill node IDs
	nodeIds := make([]cadence.Value, 0)
	for _, id := range currentEpochIdentities {
		nodeIdCdc, err := cadence.NewString(id.GetNodeID().String())
		if err != nil {
			return nil, fmt.Errorf("failed to convert node ID to cadence string %s: %w", id.GetNodeID(), err)
		}
		nodeIds = append(nodeIds, nodeIdCdc)
	}

	clusterQCAddress := systemcontracts.SystemContractsForChain(rootChainID).ClusterQC.Address.String()
	qcVoteData, err := common.ConvertClusterQcsCdc(clusterQCs, clusters, clusterQCAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to convert cluster qcs to cadence type")
	}
	currEpochFinalView, err := epoch.FinalView()
	if err != nil {
		return nil, fmt.Errorf("failed to get final view of current epoch")
	}
	currEpochTargetEndTime, err := epoch.TargetEndTime()
	if err != nil {
		return nil, fmt.Errorf("failed to get target end time of current epoch")
	}

	args := []cadence.Value{
		// recovery epoch counter
		cadence.NewUInt64(recoveryEpochCounter),
		// epoch start view
		cadence.NewUInt64(currEpochFinalView + 1),
		// staking phase end view
		cadence.NewUInt64(currEpochFinalView + numViewsInStakingAuction),
		// epoch end view
		cadence.NewUInt64(currEpochFinalView + numViewsInEpoch),
		// target duration
		cadence.NewUInt64(targetDuration),
		// target end time
		cadence.NewUInt64(currEpochTargetEndTime),
		// clusters,
		common.ConvertClusterAssignmentsCdc(assignments),
		// qcVoteData
		cadence.NewArray(qcVoteData),
		// dkg pub keys
		cadence.NewArray(dkgPubKeys),
		// dkg group key,
		dkgGroupKeyCdc,
		// dkg index map
		cadence.NewDictionary(dkgIndexMapPairs),
		// node ids
		cadence.NewArray(nodeIds),
		// recover the network by initializing a new recover epoch which will increment the smart contract epoch counter
		// or overwrite the epoch metadata for the current epoch
		cadence.NewBool(unsafeAllowOverWrite),
	}

	return args, nil
}

// ConstructRootQCsForClusters constructs a root QC for each cluster in the list.
// Args:
// - log: the logger instance.
// - clusterList: list of clusters
// - nodeInfos: list of NodeInfos (must contain all internal nodes)
// - clusterBlocks: list of root blocks (one for each cluster)
// Returns:
// - flow.AssignmentList: the generated assignment list.
// - flow.ClusterList: the generate collection cluster list.
func ConstructRootQCsForClusters(log zerolog.Logger, clusterList flow.ClusterList, nodeInfos []bootstrap.NodeInfo, clusterBlocks []*cluster.Block) []*flow.QuorumCertificate {
	if len(clusterBlocks) != len(clusterList) {
		log.Fatal().Int("len(clusterBlocks)", len(clusterBlocks)).Int("len(clusterList)", len(clusterList)).
			Msg("number of clusters needs to equal number of cluster blocks")
	}

	qcs := make([]*flow.QuorumCertificate, len(clusterBlocks))
	for i, cluster := range clusterList {
		signers := filterClusterSigners(cluster, nodeInfos)

		qc, err := GenerateClusterRootQC(signers, cluster, clusterBlocks[i])
		if err != nil {
			log.Fatal().Err(err).Int("cluster index", i).Msg("generating collector cluster root QC failed")
		}
		qcs[i] = qc
	}

	return qcs
}

// Filters a list of nodes to include only nodes that will sign the QC for the
// given cluster. The resulting list of nodes is only nodes that are in the
// given cluster AND are not partner nodes (ie. we have the private keys).
func filterClusterSigners(cluster flow.IdentitySkeletonList, nodeInfos []model.NodeInfo) []model.NodeInfo {
	var filtered []model.NodeInfo
	for _, node := range nodeInfos {
		_, isInCluster := cluster.ByNodeID(node.NodeID)
		isPrivateKeyAvailable := node.Type() == model.NodeInfoTypePrivate

		if isInCluster && isPrivateKeyAvailable {
			filtered = append(filtered, node)
		}
	}

	return filtered
}
