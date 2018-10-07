package infra

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/joeybloggs/go-download"
	"github.com/mholt/archiver"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
)

type (
	quietWriter struct{}
)

func normalDownload(ctx context.Context, uri string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		resp.Body.Close()
		return nil, fmt.Errorf("status code = %d >= 400", resp.StatusCode)
	}
	return resp.Body, nil
}

// Download is an implementation of domain.Download .
func Download(ctx context.Context, uri string, numOfDLPartitions int) (io.ReadCloser, error) {
	if numOfDLPartitions == 1 {
		return normalDownload(ctx, uri)
	}
	return download.OpenContext(
		ctx, uri, &download.Options{
			Concurrency: func(size int64) int {
				return numOfDLPartitions
			},
		})
}

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

// NewLoggerOutput returns a writer for standard logger.
func NewLoggerOutput() io.Writer {
	return quietWriter{}
}

// MkdirAll is an implementation of domain.MkdirAll .
func MkdirAll(dst string) error {
	return os.MkdirAll(dst, 0775)
}

// ReadConfigFile reads a configuration from a file.
func ReadConfigFile(dst string) (domain.Config, error) {
	cfg := domain.Config{}
	f, err := os.Open(dst)
	if err != nil {
		return cfg, err
	}
	defer f.Close()
	err = yaml.NewDecoder(f).Decode(&cfg)
	return cfg, err
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
