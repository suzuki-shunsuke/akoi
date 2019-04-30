package handler

import (
	"github.com/suzuki-shunsuke/go-cliutil"
	"github.com/urfave/cli"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
	"github.com/suzuki-shunsuke/akoi/internal/infra"
	"github.com/suzuki-shunsuke/akoi/internal/usecase/initcmd"
)

// InitCommand is the sub command "init".
var InitCommand = cli.Command{
	Name:   "init",
	Usage:  "create a configuration file if it doesn't exist",
	Action: Init,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "dest, d",
			Usage:  "created configuration file path",
			Value:  "/etc/akoi/akoi.yml",
			EnvVar: "AKOI_CONFIG_PATH",
		},
	},
}

// Init is the sub command "init".
func Init(c *cli.Context) error {
	return cliutil.ConvErrToExitError(initcmd.InitConfigFile(
		&domain.InitParams{
			Dest: c.String("dest"),
		},
		infra.FileSystem{},
	))
}
