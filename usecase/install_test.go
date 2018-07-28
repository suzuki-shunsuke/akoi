package usecase

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/suzuki-shunsuke/akoi/domain"
	"github.com/suzuki-shunsuke/akoi/registry"
	"github.com/suzuki-shunsuke/akoi/testutil"
)

func TestInstall(t *testing.T) {
	methods := &domain.InstallMethods{
		Chmod:    testutil.NewFakeChmod(nil),
		CopyFile: testutil.NewFakeCopyFile(nil),
		Download: testutil.NewFakeDownload(
			&http.Response{
				StatusCode: 200,
				Body:       testutil.NewFakeIOReadCloser("hello"),
			}, nil),
		Exist:       testutil.NewFakeExistFile(true),
		GetArchiver: testutil.NewFakeGetArchiver(nil),
		GetFileStat: testutil.NewFakeGetFileStat(
			testutil.NewFakeFileInfo("foo", 0666), nil),
		GetFileLstat: testutil.NewFakeGetFileStat(
			testutil.NewFakeFileInfo("foo", 0666), nil),
		MkdirAll: testutil.NewFakeMkdirAll(nil),
		MkLink:   testutil.NewFakeMkLink(nil),
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
		RemoveAll:  testutil.NewFakeRemoveAll(nil),
		RemoveFile: testutil.NewFakeRemoveFile(nil),
		RemoveLink: testutil.NewFakeRemoveLink(nil),
		TempDir:    testutil.NewFakeTempDir("/tmp/foo", nil),
	}
	params := &domain.InstallParams{
		ConfigFilePath: "/etc/akoi/akoi.yml", Format: "ansible"}
	if result := Install(params, methods); result.Failed {
		t.Fatal(result.String(params))
	}
	methods.ReadConfigFile = testutil.NewFakeReadConfigFile(nil, fmt.Errorf("failed to read config"))
	if result := Install(params, methods); !result.Failed {
		t.Fatal("it should be failed to read config")
	}
}

func Test_createLink(t *testing.T) {
}

func Test_installFile(t *testing.T) {
}

func Test_installPackage(t *testing.T) {
}

func Test_setupConfig(t *testing.T) {
	cfg := &domain.Config{
		BinPath:  "/usr/local/bin/{{.Name}}-{{.Version}}",
		LinkPath: "/usr/local/bin/{{.Name}}",
		Packages: map[string]domain.Package{
			"consul": {
				RawURL:  "https://releases.hashicorp.com/consul/{{.Version}}/consul_{{.Version}}_darwin_amd64.zip",
				Version: "1.2.1",
				Files: []domain.File{
					{
						Name:    "consul",
						Archive: "consul",
					},
				},
			},
		},
	}
	if err := setupConfig(cfg, registry.NewInstallMethods(true)); err != nil {
		t.Fatal(err)
	}
}
