package infra

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"

	"github.com/mholt/archiver"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
)

type (
	quietWriter struct{}
	// GetArchiver implements domain.GetArchiver .
	GetArchiver struct{}
	// GetGzipReader implements domain.GetGzipReader .
	GetGzipReader struct{}
)

// ExistFile is an implementation of domain.ExistFile .
func ExistFile(dst string) bool {
	_, err := os.Stat(dst)
	return err == nil
}

// Get converts archiver.Archiver into domain.Archiver .
func (getArchiver GetArchiver) Get(fpath, ftype string) domain.Archiver {
	if ftype == "" {
		return archiver.MatchingFormat(fpath)
	}
	arc, ok := archiver.SupportedFormats[ftype]
	if ok {
		return arc
	}
	return nil
}

// Get converts gzip.NewReader into domain.GetGzipReader .
func (getGzipReader GetGzipReader) Get(reader io.Reader) (io.ReadCloser, error) {
	return gzip.NewReader(reader)
}

// NewLoggerOutput returns a writer for standard logger.
func NewLoggerOutput() io.Writer {
	return quietWriter{}
}

// MkdirAll is an implementation of domain.MkdirAll .
func MkdirAll(dst string) error {
	return os.MkdirAll(dst, 0775)
}

// TempDir creates a temrapory directory.
func TempDir() (string, error) {
	return ioutil.TempDir("", "")
}

func (writer quietWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

// WriteFile is an implementation of domain.WriteFile .
func WriteFile(dst string, data []byte) error {
	return ioutil.WriteFile(dst, data, 0644)
}
