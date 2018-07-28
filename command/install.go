package command

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/suzuki-shunsuke/akoi/domain"
	"github.com/suzuki-shunsuke/akoi/registry"
	"github.com/suzuki-shunsuke/akoi/usecase"
)

// InstallCommand is the sub command "install".
var InstallCommand = cli.Command{
	Name:   "install",
	Usage:  "intall binaries",
	Action: Install,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "configuration file path",
			Value: "/etc/akoi/akoi.yml",
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
	params := &domain.InstallParams{
		ConfigFilePath: c.String("config"),
		Format:         c.String("format"),
		DryRun:         c.Bool("dry-run"),
	}
	result, err := usecase.Install(params, registry.NewInstallMethods(params.DryRun))
	if result == nil {
		result = &domain.Result{}
	}
	if err == nil {
		s := result.String(params)
		if s != "" {
			fmt.Println(s)
		}
		return nil
	}
	if result.Msg == "" {
		result.Msg = err.Error()
	}
	return cli.NewExitError(result.String(params), 1)
}
