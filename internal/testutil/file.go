package testutil

import (
	"context"
	"io"
	"os"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
)

// NewFakeChmod is a fake of domain.Chmod .
func NewFakeChmod(err error) domain.Chmod {
	return func(name string, mode os.FileMode) error {
		return err
	}
}

// NewFakeCopy is a fake of domain.Copy .
func NewFakeCopy(written int64, err error) domain.Copy {
	return func(dst io.Writer, src io.Reader) (int64, error) {
		return written, err
	}
}

// NewFakeDownload is a fake of domain.Download .
func NewFakeDownload(body io.ReadCloser, err error) domain.Download {
	return func(context.Context, string) (io.ReadCloser, error) {
		return body, err
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

// NewFakeNewGzipReader is a fake of domain.NewGzipReader .
func NewFakeNewGzipReader(reader io.ReadCloser, err error) domain.NewGzipReader {
	return func(r io.Reader) (io.ReadCloser, error) {
		return reader, err
	}
}

// NewFakeOpen is a fake of domain.Open .
func NewFakeOpen(f *os.File, e error) domain.Open {
	return func(name string) (*os.File, error) {
		return f, e
	}
}

// NewFakeOpenFile is a fake of domain.OpenFile .
func NewFakeOpenFile(f *os.File, e error) domain.OpenFile {
	return func(name string, flag int, perm os.FileMode) (*os.File, error) {
		return f, e
	}
}

// NewFakeReadConfigFile is a fake of domain.ReadConfigFile .
func NewFakeReadConfigFile(cfg domain.Config, err error) domain.ReadConfigFile {
	return func(dest string) (domain.Config, error) {
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
