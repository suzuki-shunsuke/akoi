package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/suzuki-shunsuke/gomic/gomic"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
	"github.com/suzuki-shunsuke/akoi/internal/test"
)

func TestInstall(t *testing.T) {
	params := domain.InstallParams{
		ConfigFilePath: "/etc/akoi/akoi.yml", Format: "ansible"}
	cfgReader := test.NewConfigReader(t, gomic.DoNothing)
	downloader := test.NewDownloader(t, gomic.DoNothing)
	getArchiver := test.NewGetArchiver(t, gomic.DoNothing)
	getGzipReader := test.NewGetGzipReader(t, gomic.DoNothing)
	result := Install(
		context.Background(), params, test.NewFileSystem(t, gomic.DoNothing),
		test.NewPrinter(t, gomic.DoNothing), cfgReader, getArchiver, downloader, getGzipReader)
	if result.Failed() {
		t.Fatal(result.String("ansible"))
	}
	cfgReader.SetReturnRead(domain.Config{}, fmt.Errorf("failed to read config"))
	result = Install(
		context.Background(), params, test.NewFileSystem(t, gomic.DoNothing),
		test.NewPrinter(t, gomic.DoNothing), cfgReader, getArchiver, downloader, getGzipReader)
	if !result.Failed() {
		t.Fatal("it should be failed to read config")
	}
}
