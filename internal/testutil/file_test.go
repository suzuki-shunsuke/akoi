package testutil

import (
	"fmt"
	"os"
	"testing"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
)

func TestFakeArchiverRead(t *testing.T) {
	arc := &FakeArchiver{err: nil}
	if err := arc.Read(nil, "foo"); err != nil {
		t.Fatal(err)
	}
}

func TestNewFakeChmod(t *testing.T) {
	f := NewFakeChmod(nil)
	if err := f("", 0755); err != nil {
		t.Fatal(err)
	}
}

func TestNewFakeCopy(t *testing.T) {
	f := NewFakeCopy(10, nil)
	if _, err := f(nil, nil); err != nil {
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
	fi := NewFakeFileInfo("foo.tar.gz", 0755)
	f := NewFakeGetFileStat(fi, nil)
	if _, err := f("foo.tar.gz"); err != nil {
		t.Fatal(err)
	}
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

func TestNewFakeOpen(t *testing.T) {
	file := os.File{}
	f := NewFakeOpen(&file, nil)
	if _, err := f("src"); err != nil {
		t.Fatal(err)
	}
}

func TestNewFakeOpenFile(t *testing.T) {
	file := os.File{}
	f := NewFakeOpenFile(&file, nil)
	if _, err := f("src", 0, 0755); err != nil {
		t.Fatal(err)
	}
}

func TestNewFakeReadConfigFile(t *testing.T) {
	cfg := domain.Config{}
	f := NewFakeReadConfigFile(cfg, nil)
	if _, err := f("src"); err != nil {
		t.Fatal(err)
	}
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
