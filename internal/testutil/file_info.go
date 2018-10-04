package testutil

import (
	"os"
	"time"
)

type (
	// FakeFileInfo implements os.FileInfo .
	FakeFileInfo struct {
		name    string
		isDir   bool
		mode    os.FileMode
		size    int64
		modTime time.Time
		sys     interface{}
	}
)

// NewFakeFileInfo creates a FakeFileInfo .
func NewFakeFileInfo(name string, mode os.FileMode) *FakeFileInfo {
	return &FakeFileInfo{
		name: name,
		mode: mode,
	}
}

// Name returns a name.
func (fi *FakeFileInfo) Name() string {
	return fi.name
}

// Size returns a size.
func (fi *FakeFileInfo) Size() int64 {
	return fi.size
}

// Mode returns a mode.
func (fi *FakeFileInfo) Mode() os.FileMode {
	return fi.mode
}

// ModTime returns a modification time.
func (fi *FakeFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir returns whether the file is a directory.
func (fi *FakeFileInfo) IsDir() bool {
	return fi.isDir
}

// Sys implements os.FileInfo#Sys() .
func (fi *FakeFileInfo) Sys() interface{} {
	return fi.sys
}
