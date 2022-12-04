package main

import (
	"fmt"
	"sync"
)

const defaultBucketCount = 256

// HashFunc is responsible for generating unsigned 64-bit hash of provided string
type HashFunc interface {
	Sum64(string) uint64
}

func NewDefaultHashFunc() HashFunc {
	return nil
}

type segment struct {
	hashmap map[uint64]uint32
}

type cache struct {
	// hashFunc represents used hash func
	hashFunc HashFunc
	// bucketCount represents the number of segments within a cache instance. value must be a power of two.
	bucketCount uint64
	// bucketMask is bitwise AND applied to the hashVal to find the segment id.
	bucketMask uint64
	// segment is shard
	segments []*segment
	// segment lock
	locks    []sync.RWMutex
	// close cache
	close chan struct{}
}

type Opt func(options *cache)


func NewCache(opts ...Opt) {
	c := &cache{
		hashFunc: NewDefaultHashFunc(),
		bucketCount: defaultBucketCount,
		close: make(chan struct{}),
	}
	for _, each := range opts {
		each(c)
	}
}

func SetShardCount(count uint64) Opt {
	return func(opt *cache) {
		opt.bucketCount = count
	}
}

func main() {
	NewCache(SetShardCount(256))
}
