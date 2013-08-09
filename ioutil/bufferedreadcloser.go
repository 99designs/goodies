package ioutil

import (
	"bytes"
	"fmt"
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
	this.ReadAll() // finish reading before we close
	return this.rc.Close()
}

func (this BufferedReadCloser) ReadAll() []byte {
	_, err := ioutil.ReadAll(this.Reader)
	fmt.Println(err)

	return this.buffer.Bytes()
}
