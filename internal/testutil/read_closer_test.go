package testutil

import (
	"testing"
)

func TestNewFakeIOReadCloser(t *testing.T) {
	if rc := NewFakeIOReadCloser(""); rc == nil {
		t.Fatal("read closer is nil")
	}
}

func TestFakeIOReadCloserRead(t *testing.T) {
	rc := NewFakeIOReadCloser("foo")
	p := []byte{}
	if _, err := rc.Read(p); err != nil {
		t.Fatal(err)
	}
}

func TestFakeIOReadCloserClose(t *testing.T) {
	rc := NewFakeIOReadCloser("foo")
	if err := rc.Close(); err != nil {
		t.Fatal(err)
	}
}
