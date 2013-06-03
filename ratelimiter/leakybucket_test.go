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
