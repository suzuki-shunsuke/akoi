package usecase

import (
	"testing"

	"github.com/suzuki-shunsuke/gomic/gomic"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
	"github.com/suzuki-shunsuke/akoi/internal/test"
)

func Test_logicSetupConfig(t *testing.T) {
	cfg := domain.Config{
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
	getArchiver := test.NewGetArchiver(t, gomic.DoNothing)
	logic := newLogicMock(t)
	logic.Fsys = test.NewFileSystem(t, gomic.DoNothing).
		SetFuncExpandEnv(func(p string) string {
			return p
		})
	if _, err := logic.SetupConfig(cfg, getArchiver); err != nil {
		t.Fatal(err)
	}
}
