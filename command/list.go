package command

import (
	"github.com/urfave/cli"

	"github.com/suzuki-shunsuke/akoi/domain"
	"github.com/suzuki-shunsuke/akoi/registry"
	"github.com/suzuki-shunsuke/akoi/usecase"
)

// ListCommand is the sub command "list".
var ListCommand = cli.Command{
	Name:   "list",
	Usage:  "list installed binaries",
	Action: List,
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
	},
}

// List is the sub command "list".
func List(c *cli.Context) error {
	params := &domain.ListParams{
		ConfigFilePath: c.String("config"),
		Format:         c.String("format"),
	}
	usecase.List(params, registry.NewListMethods(params))
	return nil
}
