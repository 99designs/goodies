package ioutil

import (
	"bytes"
	"io"
	"io/ioutil"
)

type BufferedReadCloser struct {
	io.Reader
	rc     io.ReadCloser
	buffer *bytes.Buffer
}

func NewBufferedReadCloser(r io.ReadCloser) BufferedReadCloser {
	buffer := BufferedReadCloser{
		rc:     r,
		buffer: &bytes.Buffer{},
	}
	buffer.Reader = io.TeeReader(buffer.rc, buffer.buffer)

	return buffer
}

func (this BufferedReadCloser) Close() error {
	ioutil.ReadAll(this.Reader) // finish reading everything before we close
	return this.rc.Close()
}

func (this BufferedReadCloser) ReadAll() []byte {
	ioutil.ReadAll(this.Reader)
	return this.buffer.Bytes()
}
