package storage

import "github.com/onflow/flow-go/model/flow"

type LightTransactionResultsReader interface {
	// ByBlockIDTransactionID returns the transaction result for the given block ID and transaction ID
	ByBlockIDTransactionID(blockID flow.Identifier, transactionID flow.Identifier) (*flow.LightTransactionResult, error)

	// ByBlockIDTransactionIndex returns the transaction result for the given blockID and transaction index
	ByBlockIDTransactionIndex(blockID flow.Identifier, txIndex uint32) (*flow.LightTransactionResult, error)

	// ByBlockID gets all transaction results for a block, ordered by transaction index
	ByBlockID(id flow.Identifier) ([]flow.LightTransactionResult, error)
}

// LightTransactionResults represents persistent storage for light transaction result
type LightTransactionResults interface {
	LightTransactionResultsReader

	// BatchStore inserts a batch of transaction result into a batch
	BatchStore(blockID flow.Identifier, transactionResults []flow.LightTransactionResult, rw ReaderBatchWriter) error

	// deprecated
	BatchStoreBadger(blockID flow.Identifier, transactionResults []flow.LightTransactionResult, batch BatchStorage) error
}
