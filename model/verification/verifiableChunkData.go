package verification

import (
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/state/protocol"
)

// VerifiableChunkData represents a ready-to-verify chunk
// It contains the execution result as well as all resources needed to verify it
type VerifiableChunkData struct {
	IsSystemChunk     bool                             // indicates whether this is a system chunk
	Chunk             *flow.Chunk                      // the chunk to be verified
	Header            *flow.Header                     // BlockHeader that contains this chunk
	Snapshot          protocol.SnapshotExecutionSubset // state snapshot at the chunk's block
	Result            *flow.ExecutionResult            // execution result of this block
	ChunkDataPack     *flow.ChunkDataPack              // chunk data package needed to verify this chunk
	EndState          flow.StateCommitment             // state commitment at the end of this chunk
	TransactionOffset uint32                           // index of the first transaction in a chunk within a block
}
