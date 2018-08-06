package infra

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/mholt/archiver"

	"github.com/suzuki-shunsuke/akoi/domain"
)

// ExistFile is an implementation of domain.ExistFile .
func ExistFile(dst string) bool {
	_, err := os.Stat(dst)
	return err == nil
}

// GetArchiver converts archiver.Archiver into domain.Archiver .
func GetArchiver(fpath, ftype string) domain.Archiver {
	if ftype == "" {
		return archiver.MatchingFormat(fpath)
	}
	arc, ok := archiver.SupportedFormats[ftype]
	if ok {
		return arc
	}
	return nil
}

// NewGzipReader converts gzip.NewReader into domain.NewGzipReader .
func NewGzipReader(reader io.Reader) (io.ReadCloser, error) {
	return gzip.NewReader(reader)
}

// MkdirAll is an implementation of domain.MkdirAll .
func MkdirAll(dst string) error {
	return os.MkdirAll(dst, 0775)
}

// ReadConfigFile reads a configuration from a file.
func ReadConfigFile(dst string) (*domain.Config, error) {
	f, err := os.Open(dst)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	cfg := domain.Config{}
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// TempDir creates a temrapory directory.
func TempDir() (string, error) {
	return ioutil.TempDir("", "")
}

// WriteFile is an implementation of domain.WriteFile .
func WriteFile(dst string, data []byte) error {
	return ioutil.WriteFile(dst, data, 0644)
}
