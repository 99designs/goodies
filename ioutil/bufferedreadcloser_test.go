package ioutil

import (
	"bytes"
	"io/ioutil"
	"testing"
)

const TESTIN = "test"

func new() BufferedReadCloser {
	return NewBufferedReadCloser(ioutil.NopCloser(bytes.NewReader([]byte(TESTIN))))
}

func TestReadAll(t *testing.T) {
	r := new()

	actual := string(r.ReadAll())
	if actual != TESTIN {
		t.Errorf("Expected %s, got %s", TESTIN, actual)
	}
}

func TestReadAllAfterRead(t *testing.T) {
	r := new()

	var buf []byte

	r.Read(buf)

	actual := string(r.ReadAll())
	if actual != TESTIN {
		t.Errorf("Expected %s, got %s", TESTIN, actual)
	}
}

func TestReadAllAfterClose(t *testing.T) {
	r := new()
	r.Close()

	actual := string(r.ReadAll())
	if actual != TESTIN {
		t.Errorf("Expected %s, got %s", TESTIN, actual)
	}
}

func TestReadAllAfterReadAndClose(t *testing.T) {
	r := new()

	var buf []byte

	r.Read(buf)
	r.Close()

	actual := string(r.ReadAll())
	if actual != TESTIN {
		t.Errorf("Expected %s, got %s", TESTIN, actual)
	}
}
