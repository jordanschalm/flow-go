package kvstore

import "errors"

// ErrKeyNotSet is a sentinel returned when a key is queried and no value has been set.
// The key must exist in the currently active key-value store version. This sentinel
// is used to communicate an empty/unset value rather than using zero or nil values.
// This sentinel is applicable on a key-by-key basis: some keys will always have a value
// set, others will support unset values.
var ErrKeyNotSet = errors.New("no value for requested key in Protocol State's kvstore")

// ErrKeyNotSupported is a sentinel returned when a key is read or written, but
// the key does not exist in the currently active version of the key-value store.
// This can happen in two circumstances, for example:
//  1. Current model is v2, software supports v3, and we query a key which was newly added in v3.
//  2. Current model is v3 and we query a key which was added in v2 then removed  in v3
var ErrKeyNotSupported = errors.New("protocol state's kvstore does not support the specified key at this version")

// ErrUnsupportedVersion is a sentinel returned when we attempt to decode a key-value
// store instance, but provide an unsupported version. This could happen if we accept
// an already-encoded key-value store instance from an external source (should be
// avoided in general) or if the node software version is downgraded.
var ErrUnsupportedVersion = errors.New("unsupported version for the Protocol State's kvstore")

// ErrInvalidUpgradeVersion is a sentinel returned when we attempt to set a new kvstore version
// via a ProtocolStateVersionUpgrade event, but the new version is not strictly greater than
// the current version. This error happens when smart contract has different understanding of
// the protocol state version than the node software.
var ErrInvalidUpgradeVersion = errors.New("invalid upgrade version for the Protocol State's kvstore")

// ErrInvalidActivationView is a sentinel returned when we attempt to process a KV store update,
// which has an activation view `V` so that `CurrentView + SafetyBuffer < V` does NOT hold.
var ErrInvalidActivationView = errors.New("invalid activation view for the new Protocol State version")

// ErrIncompatibleVersionChange is a sentinel returned when we attempt to replicate a parent KV store snapshot into a snapshot
// with the specified `protocolVersion` but such operation is not supported by the parent snapshot.
var ErrIncompatibleVersionChange = errors.New("incompatible version change when replicating the Protocol State's kvstore")

// ErrInvalidValue is a sentinel returned when a value is not considered valid for a given key.
// This sentinel is applicable on a key-by-key basis: same value can be considered valid/invalid for different keys.
var ErrInvalidValue = errors.New("invalid value for the requested key in Protocol State's kvstore")
