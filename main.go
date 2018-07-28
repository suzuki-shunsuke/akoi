package main

import (
	"os"

	"github.com/urfave/cli"

	"github.com/suzuki-shunsuke/akoi/command"
	"github.com/suzuki-shunsuke/akoi/domain"
)

func main() {
	app := cli.NewApp()
	app.Name = "akoi"
	app.Version = domain.Version
	app.Author = "suzuki-shunsuke https://github.com/suzuki-shunsuke"
	app.Usage = "binary version control system"
	app.Commands = []cli.Command{
		command.InitCommand,
		command.InstallCommand,
	}
	app.Run(os.Args)
}
