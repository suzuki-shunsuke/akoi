package testutil

import (
	"testing"
)

func TestNewFakeFileInfo(t *testing.T) {
	fi := NewFakeFileInfo("foo.txt", 0755)
	if fi == nil {
		t.Fatal("FileInfo is nil")
	}
}

func TestFakeFileInfoName(t *testing.T) {
	exp := "foo.txt"
	fi := NewFakeFileInfo(exp, 0755)
	act := fi.Name()
	if act != exp {
		t.Fatalf(`fi.Name() = "%s", wanted "%s"`, act, exp)
	}
}

func TestFakeFileInfoSize(t *testing.T) {
}

func TestFakeFileInfoMode(t *testing.T) {
}

func TestFakeFileInfoModTime(t *testing.T) {
}

func TestFakeFileInfoIsDir(t *testing.T) {
}

func TestFakeFileInfoSys(t *testing.T) {
}
