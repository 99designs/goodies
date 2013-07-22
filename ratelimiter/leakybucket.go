// Package ratelimiter implements the Leaky Bucket ratelimiting algorithm with memcached and in-memory backends.
package ratelimiter

import (
	"time"
)

type LeakyBucket struct {
	Size          uint16
	Fill          float64
	Leak_interval time.Duration // time.Duration for 1 unit of size to leak
	Lastupdate    time.Time
	Now           func() time.Time
}

func NewLeakyBucket(size uint16, leak_rate time.Duration) *LeakyBucket {
	bucket := LeakyBucket{
		Size:          size,
		Fill:          0,
		Leak_interval: leak_rate,
		Now:           time.Now,
		Lastupdate:    time.Now(),
	}

	return &bucket
}

func (b *LeakyBucket) updateFill() {
	now := b.Now()
	if b.Fill > 0 {
		elapsed := now.Sub(b.Lastupdate)

		b.Fill -= float64(elapsed) / float64(b.Leak_interval)
		if b.Fill < 0 {
			b.Fill = 0
		}
	}
	b.Lastupdate = now
}

func (b *LeakyBucket) Pour(amount uint16) bool {
	b.updateFill()

	var newfill float64 = b.Fill + float64(amount)

	if newfill > float64(b.Size) {
		return false
	}

	b.Fill = newfill

	return true
}

// The time at which this bucket will be completely drained
func (b *LeakyBucket) DrainedAt() time.Time {
	return b.Lastupdate.Add(b.TimeToDrain())
}

// The duration until this bucket is completely drained
func (b *LeakyBucket) TimeToDrain() time.Duration {
	return time.Duration(b.Fill / float64(b.Leak_interval))
}

type LeakyBucketSer struct {
	Size          uint16
	Fill          float64
	Leak_interval time.Duration // time.Duration for 1 unit of size to leak
	Lastupdate    time.Time
}

func (b *LeakyBucket) Serialise() *LeakyBucketSer {
	bucket := LeakyBucketSer{
		Size:          b.Size,
		Fill:          b.Fill,
		Leak_interval: b.Leak_interval,
		Lastupdate:    b.Lastupdate,
	}

	return &bucket
}

func (b *LeakyBucketSer) DeSerialise() *LeakyBucket {
	bucket := LeakyBucket{
		Size:          b.Size,
		Fill:          b.Fill,
		Leak_interval: b.Leak_interval,
		Lastupdate:    b.Lastupdate,
		Now:           time.Now,
	}

	return &bucket
}
