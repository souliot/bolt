package server

import (
	"sync/atomic"
)

// consistentIndex represents the offset of an entry in a consistent replica log.
// It implements the mvcc.ConsistentIndexGetter interface.
// It is always set to the offset of current entry before executing the entry,
// so ConsistentWatchableKV could get the consistent index from it.
type consistentIndex uint64

func (i *consistentIndex) setConsistentIndex(v uint64) {
	atomic.StoreUint64((*uint64)(i), v)
}

func (i *consistentIndex) ConsistentIndex() uint64 {
	return atomic.LoadUint64((*uint64)(i))
}
