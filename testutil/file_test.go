package testutil

import (
	"fmt"
	"testing"
)

func TestFakeArchiverRead(t *testing.T) {
}

func TestNewFakeChmod(t *testing.T) {
	f := NewFakeChmod(nil)
	if err := f("", 0755); err != nil {
		t.Fatal(err)
	}
}

func TestNewFakeDownload(t *testing.T) {
	f := NewFakeDownload(nil, nil)
	if _, err := f("http://example.com"); err != nil {
		t.Fatal(err)
	}
}

func TestNewFakeExistFile(t *testing.T) {
	f := NewFakeExistFile(true)
	if !f("") {
		t.Fatal("result must be true")
	}
	f = NewFakeExistFile(false)
	if f("") {
		t.Fatal("result must be false")
	}
}

func TestNewFakeGetArchiver(t *testing.T) {
	f := NewFakeGetArchiver(nil)
	arc := f("src.tar.gz", "")
	if arc == nil {
		t.Fatal("archiver is nil")
	}
}

func TestNewFakeGetFileStat(t *testing.T) {
}

func TestNewFakeMkdirAll(t *testing.T) {
	f := NewFakeMkdirAll(nil)
	if err := f(""); err != nil {
		t.Fatal(err)
	}
	f = NewFakeMkdirAll(fmt.Errorf("failed to create a directory"))
	if err := f(""); err == nil {
		t.Fatal("it must be failed to create a directory")
	}
}

func TestNewFakeMkLink(t *testing.T) {
	f := NewFakeMkLink(nil)
	if err := f("src", "dest"); err != nil {
		t.Fatal(err)
	}
}

func TestNewFakeReadConfigFile(t *testing.T) {
}

func TestNewFakeReadLink(t *testing.T) {
	f := NewFakeReadLink("dest", nil)
	if _, err := f("src"); err != nil {
		t.Fatal(err)
	}
}

func TestNewFakeRemoveFile(t *testing.T) {
	f := NewFakeRemoveFile(nil)
	if err := f("dest"); err != nil {
		t.Fatal(err)
	}
}

func TestNewFakeTempDir(t *testing.T) {
	f := NewFakeTempDir("dest", nil)
	if _, err := f(); err != nil {
		t.Fatal(err)
	}
}

func TestNewFakeWrite(t *testing.T) {
	f := NewFakeWrite(nil)
	if err := f("", nil); err != nil {
		t.Fatal(err)
	}
	f = NewFakeWrite(fmt.Errorf("error"))
	if err := f("", nil); err == nil {
		t.Fatal("err must not be nil")
	}
}
