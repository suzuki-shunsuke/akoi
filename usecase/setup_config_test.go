package usecase

import (
	"testing"

	"github.com/suzuki-shunsuke/akoi/domain"
	"github.com/suzuki-shunsuke/akoi/registry"
)

func Test_setupConfig(t *testing.T) {
	cfg := &domain.Config{
		BinDirTplStr:  "/usr/local/bin",
		LinkDirTplStr: "/usr/local/bin",
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
	params := domain.InstallParams{
		DryRun: true,
	}
	methods := registry.NewInstallMethods(&params)
	if err := setupConfig(cfg, &domain.SetupConfigMethods{
		GetArchiver: methods.GetArchiver,
	}); err != nil {
		t.Fatal(err)
	}
}
