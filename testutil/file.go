package testutil

import (
	"io"
	"net/http"
	"os"

	"github.com/suzuki-shunsuke/akoi/domain"
)

// NewFakeChmod is a fake of domain.Chmod .
func NewFakeChmod(err error) domain.Chmod {
	return func(name string, mode os.FileMode) error {
		return err
	}
}

// NewFakeCopyFile is a fake of domain.CopyFile .
func NewFakeCopyFile(err error) domain.CopyFile {
	return func(src, dest string) error {
		return err
	}
}

// NewFakeDownload is a fake of domain.Download .
func NewFakeDownload(resp *http.Response, err error) domain.Download {
	return func(string) (*http.Response, error) {
		return resp, err
	}
}

// NewFakeExistFile is a fake of domain.ExistFile .
func NewFakeExistFile(result bool) domain.ExistFile {
	return func(string) bool {
		return result
	}
}

// NewFakeGetArchiver is a fake of domain.GetArchiver .
func NewFakeGetArchiver(err error) domain.GetArchiver {
	return func(fpath, ftype string) domain.Archiver {
		return &FakeArchiver{err: err}
	}
}

// NewFakeGetFileStat is a fake of domain.GetFileStat .
func NewFakeGetFileStat(fi os.FileInfo, err error) domain.GetFileStat {
	return func(string) (os.FileInfo, error) {
		return fi, err
	}
}

// NewFakeMkdirAll is a fake of domain.MkdirAll .
func NewFakeMkdirAll(e error) domain.MkdirAll {
	return func(string) error {
		return e
	}
}

// NewFakeMkLink is a fake of domain.MkLink .
func NewFakeMkLink(e error) domain.MkLink {
	return func(src, dest string) error {
		return e
	}
}

// NewFakeReadConfigFile is a fake of domain.ReadConfigFile .
func NewFakeReadConfigFile(cfg *domain.Config, err error) domain.ReadConfigFile {
	return func(dest string) (*domain.Config, error) {
		return cfg, err
	}
}

// NewFakeReadLink is a fake of domain.ReadLink .
func NewFakeReadLink(dest string, err error) domain.ReadLink {
	return func(src string) (string, error) {
		return dest, err
	}
}

// NewFakeRemoveFile is a fake of domain.RemoveFile .
func NewFakeRemoveFile(e error) domain.RemoveFile {
	return func(dest string) error {
		return e
	}
}

// NewFakeTempDir is a fake of domain.TempDir .
func NewFakeTempDir(dst string, err error) domain.TempDir {
	return func() (string, error) {
		return dst, err
	}
}

// NewFakeWrite is a fake of domain.WriteFile .
func NewFakeWrite(e error) domain.WriteFile {
	return func(dest string, data []byte) error {
		return e
	}
}

// FakeArchiver is a fake of domain.Archiver .
type FakeArchiver struct {
	err error
}

// Read implements domain.Archiver#Read .
func (arc *FakeArchiver) Read(input io.Reader, destination string) error {
	return arc.err
}
