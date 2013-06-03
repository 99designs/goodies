// Package expect has a few test utils.
package expect

import (
	"bytes"
	"testing"
	"time"
)

func Expect(t *testing.T, name string, actual interface{}, expected interface{}) {
	if actual != expected {
		t.Errorf("Expected %s to be %v, got %v", name, expected, actual)
	}
}

func ExpectBinary(t *testing.T, name string, actual string, expected string) {
	if actual != expected {
		t.Errorf("Expected %s to be \n%q, got \n%q", name, expected, actual)
	}
}

func ExpectTime(t *testing.T, name string, actual time.Time, expected time.Time) {
	if actual != expected {
		t.Errorf("Expected %s to be \n%d, got \n%d", name, expected.Unix(), actual.Unix())
	}
}

func ExpectBytes(t *testing.T, name string, actual, expected []byte) {
	if !bytes.Equal(expected, actual) {
		t.Errorf("Expected %s to be \n%q, got \n%q", name, expected, actual)
	}
}
