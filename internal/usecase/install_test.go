package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/suzuki-shunsuke/gomic/gomic"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
	"github.com/suzuki-shunsuke/akoi/internal/test"
)

func Test_logicInstall(t *testing.T) {
	params := domain.InstallParams{
		ConfigFilePath: "/etc/akoi/akoi.yml", Format: "ansible"}
	cfgReader := test.NewConfigReader(t, gomic.DoNothing)
	downloader := test.NewDownloader(t, gomic.DoNothing)
	getArchiver := test.NewGetArchiver(t, gomic.DoNothing)
	getGzipReader := test.NewGetGzipReader(t, gomic.DoNothing)
	lgc := newLogicMock(t)
	result := lgc.Install(
		context.Background(), params,
		test.NewPrinter(t, gomic.DoNothing), cfgReader, getArchiver, downloader, getGzipReader)
	if result.Failed() {
		t.Fatal(result.String("ansible"))
	}
	cfgReader.SetReturnRead(domain.Config{}, fmt.Errorf("failed to read config"))
	result = lgc.Install(
		context.Background(), params,
		test.NewPrinter(t, gomic.DoNothing), cfgReader, getArchiver, downloader, getGzipReader)
	if !result.Failed() {
		t.Fatal("it should be failed to read config")
	}
}
