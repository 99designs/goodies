package ratelimiter

import (
	"testing"
	"time"
)

func TestPour(t *testing.T) {
	bucket := NewLeakyBucket(60, time.Second)
	bucket.Lastupdate = time.Unix(0, 0)

	bucket.Now = func() time.Time { return time.Unix(1, 0) }

	if bucket.Pour(61) {
		t.Error("Expected false")
	}

	if !bucket.Pour(10) {
		t.Error("Expected true")
	}

	if !bucket.Pour(49) {
		t.Error("Expected true")
	}

	if bucket.Pour(2) {
		t.Error("Expected false")
	}

	bucket.Now = func() time.Time { return time.Unix(61, 0) }
	if !bucket.Pour(60) {
		t.Error("Expected true")
	}

	if bucket.Pour(1) {
		t.Error("Expected false")
	}

	bucket.Now = func() time.Time { return time.Unix(70, 0) }

	if !bucket.Pour(1) {
		t.Error("Expected true")
	}

}

func TestTimeSinceLastUpdate(t *testing.T) {
	bucket := NewLeakyBucket(60, time.Second)
	bucket.Pour(1)
	bucket.Lastupdate = bucket.Lastupdate.Add(-time.Second)
	sinceLast := bucket.TimeSinceLastUpdate()
	if sinceLast < time.Second {
		t.Error("Expected time since last update to be less than 1 second, got %+v", sinceLast)
	}
	if sinceLast > (time.Millisecond * 1100) {
		t.Error("Expected time since last update to be about 1 second, got %+v", sinceLast)
	}
}

func TestTimeToDrain(t *testing.T) {
	bucket := NewLeakyBucket(60, time.Second)
	bucket.Pour(10)
	ttd := bucket.TimeToDrain()
	if ttd > time.Second*10 {
		t.Error("Time to drain should be <= 10 seconds")
	}
	if ttd <= time.Second*9 {
		t.Error("Time to drain should be > 9 seconds")
	}
}
