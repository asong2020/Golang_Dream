package main

import (
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

type options struct {
	hashFunc HashFunc
	bucketCount uint64
}

type Option interface {
	apply(*options)
}

type Bucket struct {
	count uint64
}

func (b Bucket) apply(opts *options) {
	opts.bucketCount = b.count
}

func WithBucketCount(count uint64) Option {
	return Bucket{
		count: count,
	}
}

type Hash struct {
	hashFunc HashFunc
}

func (h Hash) apply(opts *options)  {
	opts.hashFunc = h.hashFunc
}

func WithHashFunc(hashFunc HashFunc) Option {
	return Hash{hashFunc: hashFunc}
}

func NewCache(opts ...Option) {
	o := &options{
		hashFunc: NewDefaultHashFunc(),
		bucketCount: defaultBucketCount,
	}
	for _, each := range opts {
		each.apply(o)
	}
}

func main() {
	NewCache(WithBucketCount(128))
}
