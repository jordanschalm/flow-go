package operation

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/dgraph-io/badger/v2"
	"github.com/vmihailenco/msgpack/v4"

	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/module/irrecoverable"
	"github.com/onflow/flow-go/storage"
)

// batchWrite will encode the given entity using msgpack and will upsert the resulting
// binary data in the badger wrote batch under the provided key - if the value already exists
// in the database it will be overridden.
// No errors are expected during normal operation.
func batchWrite(key []byte, entity interface{}) func(writeBatch *badger.WriteBatch) error {
	return func(writeBatch *badger.WriteBatch) error {

		// update the maximum key size if the inserted key is bigger
		if uint32(len(key)) > max {
			max = uint32(len(key))
			err := SetMax(writeBatch)
			if err != nil {
				return fmt.Errorf("could not update max tracker: %w", err)
			}
		}

		// serialize the entity data
		val, err := msgpack.Marshal(entity)
		if err != nil {
			return irrecoverable.NewExceptionf("could not encode entity: %w", err)
		}

		// persist the entity data into the DB
		err = writeBatch.Set(key, val)
		if err != nil {
			return irrecoverable.NewExceptionf("could not store data: %w", err)
		}
		return nil
	}
}

// insert will encode the given entity using msgpack and will insert the resulting
// binary data in the badger DB under the provided key. It will error if the
// key already exists.
// Error returns:
//   - storage.ErrAlreadyExists if the key already exists in the database.
//   - generic error in case of unexpected failure from the database layer or
//     encoding failure.
func insert(key []byte, entity interface{}) func(*badger.Txn) error {
	return func(tx *badger.Txn) error {

		// update the maximum key size if the inserted key is bigger
		if uint32(len(key)) > max {
			max = uint32(len(key))
			err := SetMax(tx)
			if err != nil {
				return fmt.Errorf("could not update max tracker: %w", err)
			}
		}

		// check if the key already exists in the db
		_, err := tx.Get(key)
		if err == nil {
			return storage.ErrAlreadyExists
		}

		if !errors.Is(err, badger.ErrKeyNotFound) {
			return irrecoverable.NewExceptionf("could not retrieve key: %w", err)
		}

		// serialize the entity data
		val, err := msgpack.Marshal(entity)
		if err != nil {
			return irrecoverable.NewExceptionf("could not encode entity: %w", err)
		}

		// persist the entity data into the DB
		err = tx.Set(key, val)
		if err != nil {
			return irrecoverable.NewExceptionf("could not store data: %w", err)
		}
		return nil
	}
}

// update will encode the given entity with MsgPack and update the binary data
// under the given key in the badger DB. The key must already exist.
// Error returns:
//   - storage.ErrNotFound if the key does not already exist in the database.
//   - generic error in case of unexpected failure from the database layer or
//     encoding failure.
func update(key []byte, entity interface{}) func(*badger.Txn) error {
	return func(tx *badger.Txn) error {

		// retrieve the item from the key-value store
		_, err := tx.Get(key)
		if errors.Is(err, badger.ErrKeyNotFound) {
			return storage.ErrNotFound
		}
		if err != nil {
			return irrecoverable.NewExceptionf("could not check key: %w", err)
		}

		// serialize the entity data
		val, err := msgpack.Marshal(entity)
		if err != nil {
			return irrecoverable.NewExceptionf("could not encode entity: %w", err)
		}

		// persist the entity data into the DB
		err = tx.Set(key, val)
		if err != nil {
			return irrecoverable.NewExceptionf("could not replace data: %w", err)
		}

		return nil
	}
}

// upsert will encode the given entity with MsgPack and upsert the binary data
// under the given key in the badger DB.
func upsert(key []byte, entity interface{}) func(*badger.Txn) error {
	return func(tx *badger.Txn) error {
		// update the maximum key size if the inserted key is bigger
		if uint32(len(key)) > max {
			max = uint32(len(key))
			err := SetMax(tx)
			if err != nil {
				return fmt.Errorf("could not update max tracker: %w", err)
			}
		}

		// serialize the entity data
		val, err := msgpack.Marshal(entity)
		if err != nil {
			return irrecoverable.NewExceptionf("could not encode entity: %w", err)
		}

		// persist the entity data into the DB
		err = tx.Set(key, val)
		if err != nil {
			return irrecoverable.NewExceptionf("could not upsert data: %w", err)
		}

		return nil
	}
}

// remove removes the entity with the given key, if it exists. If it doesn't
// exist, this is a no-op.
// Error returns:
// * storage.ErrNotFound if the key to delete does not exist.
// * generic error in case of unexpected database error
func remove(key []byte) func(*badger.Txn) error {
	return func(tx *badger.Txn) error {
		// retrieve the item from the key-value store
		_, err := tx.Get(key)
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return storage.ErrNotFound
			}
			return irrecoverable.NewExceptionf("could not check key: %w", err)
		}

		err = tx.Delete(key)
		if err != nil {
			return irrecoverable.NewExceptionf("could not delete item: %w", err)
		}
		return nil
	}
}

// batchRemove removes entry under a given key in a write-batch.
// if key doesn't exist, does nothing.
// No errors are expected during normal operation.
func batchRemove(key []byte) func(writeBatch *badger.WriteBatch) error {
	return func(writeBatch *badger.WriteBatch) error {
		err := writeBatch.Delete(key)
		if err != nil {
			return irrecoverable.NewExceptionf("could not batch delete data: %w", err)
		}
		return nil
	}
}

// removeByPrefix removes all the entities if the prefix of the key matches the given prefix.
// if no key matches, this is a no-op
// No errors are expected during normal operation.
func removeByPrefix(prefix []byte) func(*badger.Txn) error {
	return func(tx *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.AllVersions = false
		opts.PrefetchValues = false
		it := tx.NewIterator(opts)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			key := it.Item().KeyCopy(nil)
			err := tx.Delete(key)
			if err != nil {
				return irrecoverable.NewExceptionf("could not delete item with prefix: %w", err)
			}
		}

		return nil
	}
}

// retrieve will retrieve the binary data under the given key from the badger DB
// and decode it into the given entity. The provided entity needs to be a
// pointer to an initialized entity of the correct type.
// Error returns:
//   - storage.ErrNotFound if the key does not exist in the database
//   - generic error in case of unexpected failure from the database layer, or failure
//     to decode an existing database value
func retrieve(key []byte, entity interface{}) func(*badger.Txn) error {
	return func(tx *badger.Txn) error {

		// retrieve the item from the key-value store
		item, err := tx.Get(key)
		if errors.Is(err, badger.ErrKeyNotFound) {
			return storage.ErrNotFound
		}
		if err != nil {
			return irrecoverable.NewExceptionf("could not load data: %w", err)
		}

		// get the value from the item
		err = item.Value(func(val []byte) error {
			err := msgpack.Unmarshal(val, entity)
			return err
		})
		if err != nil {
			return irrecoverable.NewExceptionf("could not decode entity: %w", err)
		}

		return nil
	}
}

// exists returns true if a key exists in the database.
// No errors are expected during normal operation.
func exists(key []byte, keyExists *bool) func(*badger.Txn) error {
	return func(tx *badger.Txn) error {
		_, err := tx.Get(key)
		if err != nil {
			// the key does not exist in the database
			if errors.Is(err, badger.ErrKeyNotFound) {
				*keyExists = false
				return nil
			}
			// exception while checking for the key
			return irrecoverable.NewExceptionf("could not load data: %w", err)
		}

		// the key does exist in the database
		*keyExists = true
		return nil
	}
}

// checkFunc is called during key iteration through the badger DB in order to
// check whether we should process the given key-value pair. It can be used to
// avoid loading the value if its not of interest, as well as storing the key
// for the current iteration step.
type checkFunc func(key []byte) bool

// createFunc returns a pointer to an initialized entity that we can potentially
// decode the next value into during a badger DB iteration.
type createFunc func() interface{}

// handleFunc is a function that starts the processing of the current key-value
// pair during a badger iteration. It should be called after the key was checked
// and the entity was decoded.
// No errors are expected during normal operation. Any errors will halt the iteration.
type handleFunc func() error

// iterationFunc is a function provided to our low-level iteration function that
// allows us to pass badger efficiencies across badger boundaries. By calling it
// for each iteration step, we can inject a function to check the key, a
// function to create the decode target and a function to process the current
// key-value pair. This a consumer of the API to decode when to skip the loading
// of values, the initialization of entities and the processing.
type iterationFunc func() (checkFunc, createFunc, handleFunc)

// lookup is the default iteration function allowing us to collect a list of
// entity IDs from an index.
func lookup(entityIDs *[]flow.Identifier) func() (checkFunc, createFunc, handleFunc) {
	*entityIDs = make([]flow.Identifier, 0, len(*entityIDs))
	return func() (checkFunc, createFunc, handleFunc) {
		check := func(key []byte) bool {
			return true
		}
		var entityID flow.Identifier
		create := func() interface{} {
			return &entityID
		}
		handle := func() error {
			*entityIDs = append(*entityIDs, entityID)
			return nil
		}
		return check, create, handle
	}
}

// withPrefetchValuesFalse configures a Badger iteration to NOT preemptively load
// the values when iterating over keys (ie. key-only iteration). Key-only iteration
// is several order of magnitudes faster than regular iteration, because it involves
// access to the LSM-tree only, which is usually resident entirely in RAM.
func withPrefetchValuesFalse(options *badger.IteratorOptions) {
	options.PrefetchValues = false
}

// iterate iterates over a range of keys defined by a start and end key. The
// start key may be higher than the end key, in which case we iterate in
// reverse order.
//
// The iteration range uses prefix-wise semantics. Specifically, all keys that
// meet ANY of the following conditions are included in the iteration:
//   - have a prefix equal to the start key OR
//   - have a prefix equal to the end key OR
//   - have a prefix that is lexicographically between start and end
//
// On each iteration, it will call the iteration function to initialize
// functions specific to processing the given key-value pair.
//
// TODO: this function is unbounded – pass context.Context to this or calling functions to allow timing functions out.
// No errors are expected during normal operation. Any errors returned by the
// provided handleFunc will be propagated back to the caller of iterate.
func iterate(start []byte, end []byte, iteration iterationFunc, opts ...func(*badger.IteratorOptions)) func(*badger.Txn) error {
	return func(tx *badger.Txn) error {

		// initialize the default options and comparison modifier for iteration
		modifier := 1
		options := badger.DefaultIteratorOptions
		for _, apply := range opts {
			apply(&options)
		}

		// In order to satisfy this function's prefix-wise inclusion semantics,
		// we append 0xff bytes to the largest of start and end.
		// This ensures Badger will seek to the largest key with that prefix
		// for reverse iteration, thus including all keys with a prefix matching
		// the starting key. It also enables us to detect boundary conditions by
		// simple lexicographic comparison (ie. bytes.Compare) rather than
		// explicitly comparing prefixes.
		//
		// See https://github.com/onflow/flow-go/pull/3310#issuecomment-618127494
		// for discussion and more detail on this.

		// If start is bigger than end, we have a backwards iteration:
		// 1) We set the reverse option on the iterator, so we step through all
		//    the keys backwards. This modifies the behaviour of Seek to go to
		//    the first key that is less than or equal to the start key (as
		//    opposed to greater than or equal in a regular iteration).
		// 2) In order to satisfy this function's prefix-wise inclusion semantics,
		//    we append a 0xff-byte suffix to the start key so the seek will go
		// to the right place.
		// 3) For a regular iteration, we break the loop upon hitting the first
		//    item that has a key higher than the end prefix. In order to reverse
		//    this, we use a modifier for the comparison that reverses the check
		//    and makes it stop upon the first item lower than the end prefix.
		if bytes.Compare(start, end) > 0 {
			options.Reverse = true // make sure to go in reverse order
			modifier = -1          // make sure to stop after end prefix
			length := uint32(len(start))
			diff := max - length
			for i := uint32(0); i < diff; i++ {
				start = append(start, 0xff)
			}
		} else {
			// for forward iteration, add the 0xff-bytes suffix to the end
			// prefix, to ensure we include all keys with that prefix before
			// finishing.
			length := uint32(len(end))
			diff := max - length
			for i := uint32(0); i < diff; i++ {
				end = append(end, 0xff)
			}
		}

		it := tx.NewIterator(options)
		defer it.Close()

		for it.Seek(start); it.Valid(); it.Next() {

			item := it.Item()

			key := item.Key()
			// for forward iteration, check whether key > end, for backward
			// iteration check whether key < end
			if bytes.Compare(key, end)*modifier > 0 {
				break
			}

			// initialize processing functions for iteration
			check, create, handle := iteration()

			// check if we should process the item at all
			ok := check(key)
			if !ok {
				continue
			}

			// process the actual item
			err := item.Value(func(val []byte) error {

				// decode into the entity
				entity := create()
				err := msgpack.Unmarshal(val, entity)
				if err != nil {
					return irrecoverable.NewExceptionf("could not decode entity: %w", err)
				}

				// process the entity
				err = handle()
				if err != nil {
					return fmt.Errorf("could not handle entity: %w", err)
				}

				return nil
			})
			if err != nil {
				return fmt.Errorf("could not process value: %w", err)
			}
		}

		return nil
	}
}

// traverse iterates over a range of keys defined by a prefix.
//
// The prefix must be shared by all keys in the iteration.
//
// On each iteration, it will call the iteration function to initialize
// functions specific to processing the given key-value pair.
func traverse(prefix []byte, iteration iterationFunc) func(*badger.Txn) error {
	return func(tx *badger.Txn) error {
		if len(prefix) == 0 {
			return fmt.Errorf("prefix must not be empty")
		}

		opts := badger.DefaultIteratorOptions
		// NOTE: this is an optimization only, it does not enforce that all
		// results in the iteration have this prefix.
		opts.Prefix = prefix

		it := tx.NewIterator(opts)
		defer it.Close()

		// this is where we actually enforce that all results have the prefix
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {

			item := it.Item()

			// initialize processing functions for iteration
			check, create, handle := iteration()

			// check if we should process the item at all
			key := item.Key()
			ok := check(key)
			if !ok {
				continue
			}

			// process the actual item
			err := item.Value(func(val []byte) error {

				// decode into the entity
				entity := create()
				err := msgpack.Unmarshal(val, entity)
				if err != nil {
					return irrecoverable.NewExceptionf("could not decode entity: %w", err)
				}

				// process the entity
				err = handle()
				if err != nil {
					return fmt.Errorf("could not handle entity: %w", err)
				}

				return nil
			})
			if err != nil {
				return fmt.Errorf("could not process value: %w", err)
			}
		}

		return nil
	}
}

// findHighestAtOrBelow searches for the highest key with the given prefix and a height
// at or below the target height, and retrieves and decodes the value associated with the
// key into the given entity.
// If no key is found, the function returns storage.ErrNotFound.
func findHighestAtOrBelow(
	prefix []byte,
	height uint64,
	entity interface{},
) func(*badger.Txn) error {
	return func(tx *badger.Txn) error {
		if len(prefix) == 0 {
			return fmt.Errorf("prefix must not be empty")
		}

		opts := badger.DefaultIteratorOptions
		opts.Prefix = prefix
		opts.Reverse = true

		it := tx.NewIterator(opts)
		defer it.Close()

		it.Seek(append(prefix, keyPartToBinary(height)...))

		if !it.Valid() {
			return storage.ErrNotFound
		}

		return it.Item().Value(func(val []byte) error {
			err := msgpack.Unmarshal(val, entity)
			if err != nil {
				return fmt.Errorf("could not decode entity: %w", err)
			}
			return nil
		})
	}
}

// Fail returns a DB operation function that always fails with the given error.
func Fail(err error) func(*badger.Txn) error {
	return func(_ *badger.Txn) error {
		return err
	}
}
