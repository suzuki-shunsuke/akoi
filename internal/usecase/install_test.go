package usecase

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
	"github.com/suzuki-shunsuke/akoi/internal/infra"
	"github.com/suzuki-shunsuke/akoi/internal/testutil"
)

func TestInstall(t *testing.T) {
	methods := &domain.InstallMethods{
		Chmod: testutil.NewFakeChmod(nil),
		Copy:  testutil.NewFakeCopy(10, nil),
		Download: testutil.NewFakeDownload(
			&http.Response{
				StatusCode: 200,
				Body:       testutil.NewFakeIOReadCloser("hello"),
			}, nil),
		ExpandEnv:   os.ExpandEnv,
		Fprintf:     infra.NewFprintf(true),
		Fprintln:    infra.NewFprintln(true),
		GetArchiver: testutil.NewFakeGetArchiver(nil),
		GetFileStat: testutil.NewFakeGetFileStat(
			testutil.NewFakeFileInfo("foo", 0666), nil),
		GetFileLstat: testutil.NewFakeGetFileStat(
			testutil.NewFakeFileInfo("foo", 0666), nil),
		MkdirAll: testutil.NewFakeMkdirAll(nil),
		MkLink:   testutil.NewFakeMkLink(nil),
		NewGzipReader: testutil.NewFakeNewGzipReader(
			testutil.NewFakeIOReadCloser("hello"), nil),
		Open:     testutil.NewFakeOpen(&os.File{}, nil),
		OpenFile: testutil.NewFakeOpenFile(&os.File{}, nil),
		Printf:   infra.NewPrintf(true),
		Println:  infra.NewPrintln(true),
		ReadConfigFile: testutil.NewFakeReadConfigFile(
			&domain.Config{
				BinPath:  "/usr/local/bin/{{.Name}}-{{.Version}}",
				LinkPath: "/usr/local/bin/{{.Name}}",
				Packages: map[string]domain.Package{
					"consul": {
						RawURL:  "https://releases.hashicorp.com/consul/{{.Version}}/consul_{{.Version}}_darwin_amd64.zip",
						Version: "1.2.0",
						Files: []domain.File{
							{
								Name:    "consul",
								Archive: "consul",
							},
						},
					},
				},
			}, nil),
		ReadLink:   testutil.NewFakeReadLink("/usr/local/bin/consul", nil),
		RemoveAll:  testutil.NewFakeRemoveFile(nil),
		RemoveFile: testutil.NewFakeRemoveFile(nil),
		RemoveLink: testutil.NewFakeRemoveFile(nil),
		TempDir:    testutil.NewFakeTempDir("/tmp/foo", nil),
	}
	params := &domain.InstallParams{
		ConfigFilePath: "/etc/akoi/akoi.yml", Format: "ansible"}
	if result := Install(context.Background(), params, methods); result.Failed {
		t.Fatal(result.String(params))
	}
	methods.ReadConfigFile = testutil.NewFakeReadConfigFile(
		nil, fmt.Errorf("failed to read config"))
	if result := Install(context.Background(), params, methods); !result.Failed {
		t.Fatal("it should be failed to read config")
	}
}
