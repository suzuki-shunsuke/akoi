package handler

import (
	"context"
	"fmt"

	"github.com/urfave/cli"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
	"github.com/suzuki-shunsuke/akoi/internal/registry"
	"github.com/suzuki-shunsuke/akoi/internal/usecase"
)

// InstallCommand is the sub command "install".
var InstallCommand = cli.Command{
	Name:   "install",
	Usage:  "intall binaries",
	Action: Install,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "configuration file path",
			Value:  "/etc/akoi/akoi.yml",
			EnvVar: "AKOI_CONFIG_PATH",
		},
		cli.StringFlag{
			Name:  "format, f",
			Usage: "output format",
			Value: "human",
		},
		cli.BoolFlag{
			Name:  "dry-run, n",
			Usage: "dry run",
		},
	},
}

// Install is the sub command "install".
func Install(c *cli.Context) error {
	params := domain.InstallParams{
		ConfigFilePath: c.String("config"),
		Format:         c.String("format"),
		DryRun:         c.Bool("dry-run"),
	}
	result := usecase.Install(
		context.Background(), params, registry.NewInstallMethods(params))
	if result == nil {
		result = &domain.Result{}
	}
	if !result.Failed {
		s := result.String(params.Format)
		if s != "" {
			fmt.Println(s)
		}
		return nil
	}
	return cli.NewExitError(result.String(params.Format), 1)
}
