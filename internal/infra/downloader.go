package infra

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/joeybloggs/go-download"
)

type (
	// Downloader implements domain.Downloader .
	Downloader struct{}
)

// Download downloads a file.
func (dl Downloader) Download(ctx context.Context, uri string, numOfDLPartitions int) (io.ReadCloser, error) {
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
