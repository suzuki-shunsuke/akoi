package infra

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestExistFile(t *testing.T) {
	p := "file_test.go"
	if !ExistFile(p) {
		t.Fatalf("%s must exist", p)
	}
	p = "_.go"
	if ExistFile(p) {
		t.Fatalf("%s must not exist", p)
	}
}

func TestWriteFile(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	if err := WriteFile(f.Name(), []byte("test of akoi")); err != nil {
		t.Fatal(err)
	}
}

func TestMkdirAll(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(dir)
	if err := MkdirAll(filepath.Join(dir, "foo")); err != nil {
		t.Fatal(err)
	}
}

func TestMkLink(t *testing.T) {
}

func TestReadConfigFile(t *testing.T) {
}

func TestReadLink(t *testing.T) {
}

func TestRemoveAll(t *testing.T) {
}

func TestRemoveLink(t *testing.T) {
}

func TestTempDir(t *testing.T) {
}
