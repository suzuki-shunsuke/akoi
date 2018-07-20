package testutil

import (
	"bytes"
)

type (
	// FakeIOReadCloser implements io.Reader and io.Closer .
	FakeIOReadCloser struct {
		buf *bytes.Buffer
	}
)

// NewFakeIOReadCloser returns a new FakeIOReadCloser .
func NewFakeIOReadCloser(data string) *FakeIOReadCloser {
	return &FakeIOReadCloser{
		buf: bytes.NewBufferString(data),
	}
}

// Read implements io.Reader .
func (r *FakeIOReadCloser) Read(p []byte) (int, error) {
	return r.buf.Read(p)
}

// Close implements io.Closer .
func (r *FakeIOReadCloser) Close() error {
	return nil
}
