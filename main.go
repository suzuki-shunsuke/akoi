package main

import (
	"os"

	"github.com/urfave/cli"

	"github.com/suzuki-shunsuke/akoi/command"
)

func main() {
	app := cli.NewApp()
	app.Name = "akoi"
	app.Version = "0.3.1"
	app.Author = "suzuki-shunsuke https://github.com/suzuki-shunsuke"
	app.Usage = "binary version control system"
	app.Commands = []cli.Command{
		{
			Name:   "init",
			Usage:  "create a configuration file if it doesn't exist",
			Action: command.Init,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "dest, d",
					Usage: "created configuration file path",
					Value: "/etc/akoi/akoi.yml",
				},
			},
		},
		{
			Name:   "install",
			Usage:  "intall binaries",
			Action: command.Install,
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
		},
	}
	app.Run(os.Args)
}
