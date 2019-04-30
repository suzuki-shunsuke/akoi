package infra

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/joeybloggs/go-download"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
)

type (
	// Downloader implements domain.Downloader .
	Downloader struct{}
)

// Download downloads a file.
func (dl Downloader) Download(
	ctx context.Context, uri string, option domain.DownloadOption,
) (io.ReadCloser, error) {
	if option.DLPartitionCount == 1 {
		return normalDownload(ctx, uri, option)
	}
	var clFunc download.ClientFn
	if option.Timeout != 0 {
		clFunc = func() http.Client {
			return http.Client{
				Timeout: time.Duration(option.Timeout),
			}
		}
	}
	return download.OpenContext(
		ctx, uri, &download.Options{
			Concurrency: func(size int64) int {
				return option.DLPartitionCount
			},
			Client: clFunc,
		})
}

func normalDownload(
	ctx context.Context, uri string, option domain.DownloadOption,
) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	client := http.DefaultClient
	if option.Timeout != 0 {
		client.Timeout = time.Duration(option.Timeout)
	}
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
