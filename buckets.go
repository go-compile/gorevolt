package gorevolt

import (
	"sync"
	"time"
)

const (
	bucketReplenish = 10
)

var (
	bucketUsers      = newBucket(20)
	bucketBots       = newBucket(10)
	bucketChannels   = newBucket(15)
	bucketServers    = newBucket(5)
	bucketAuth       = newBucket(3)
	bucketAuthDelete = newBucket(255)
)

// bucket is used to manage rate limiting
type bucket struct {
	s         int64
	hits, max int
	sync.Mutex
}

func newBucket(max int) *bucket {
	return &bucket{max: max}
}

func (b *bucket) Increment(n int) int {
	b.Lock()
	defer b.Unlock()

	if b.hits == 0 {
		b.s = time.Now().Unix()
	}

	// check if bucket has been replenished
	if b.s+bucketReplenish < time.Now().Unix() {
		// reset bucket
		b.hits = 0
	}

	b.hits += n
	return b.hits
}

// Wait will only wait for the bucket to refresh if its at its maximum
func (b *bucket) Wait() {
	b.Lock()

	if b.hits >= b.max {
		b.Unlock()
		time.Sleep(time.Until(time.Unix(b.s+bucketReplenish, 0)))
		return
	}

	b.Unlock()
}
