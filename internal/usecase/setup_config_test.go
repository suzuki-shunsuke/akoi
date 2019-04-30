package usecase

import (
	"testing"

	"github.com/stretchr/testify/require"
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
	logic := newLogicMock(t)
	logic.Fsys = test.NewFileSystem(t, gomic.DoNothing).
		SetFuncExpandEnv(func(p string) string {
			return p
		})

	_, err := logic.SetupConfig(cfg, ".akoi.yml")
	require.Nil(t, err)
}

func Test_logicSetupPkgConfig(t *testing.T) {
	pkg := domain.Package{
		RawURL:  "https://releases.hashicorp.com/consul/{{.Version}}/consul_{{.Version}}_darwin_amd64.zip",
		Version: "1.2.1",
		Files: []domain.File{
			{
				Name:    "consul",
				Archive: "consul",
			},
		},
	}
	pkgName := "consul"
	cfg := domain.Config{
		BinPath:  "/usr/local/bin/{{.Name}}-{{.Version}}",
		LinkPath: "/usr/local/bin/{{.Name}}",
		Packages: map[string]domain.Package{
			pkgName: pkg,
		},
	}
	logic := newLogicMock(t)
	logic.Fsys = test.NewFileSystem(t, gomic.DoNothing).
		SetFuncExpandEnv(func(p string) string {
			return p
		})
	_, err := logic.SetupPkgConfig(cfg, "/etc/akoi", pkgName, pkg, 4)
	require.Nil(t, err)
}

func Test_logicSetupFileConfig(t *testing.T) {
	file := domain.File{
		Name:     "consul",
		Archive:  "consul",
		BinPath:  "/usr/local/bin/akoi-{{.Version}}",
		LinkPath: "/usr/local/bin/akoi",
	}
	pkg := domain.Package{
		RawURL:  "https://releases.hashicorp.com/consul/{{.Version}}/consul_{{.Version}}_darwin_amd64.zip",
		Version: "1.2.1",
		Files:   []domain.File{file},
	}
	logic := newLogicMock(t)
	logic.Fsys = test.NewFileSystem(t, gomic.DoNothing).
		SetFuncExpandEnv(func(p string) string {
			return p
		})
	_, err := logic.SetupFileConfig(pkg, "/etc/akoi", file)
	require.Nil(t, err)
}
