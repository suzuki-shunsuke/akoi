package usecase

import (
	"testing"

	"github.com/suzuki-shunsuke/akoi/domain"
	"github.com/suzuki-shunsuke/akoi/registry"
)

func Test_setupConfig(t *testing.T) {
	cfg := &domain.Config{
		BinDir:  "/usr/local/bin",
		LinkDir: "/usr/local/bin",
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
	if err := setupConfig(cfg, registry.NewInstallMethods(&params)); err != nil {
		t.Fatal(err)
	}
}
