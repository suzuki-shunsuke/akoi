package main

import (
	"os"

	"github.com/urfave/cli"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
	"github.com/suzuki-shunsuke/akoi/internal/handler"
)

func main() {
	app := cli.NewApp()
	app.Name = "akoi"
	app.Version = domain.Version
	app.Author = "suzuki-shunsuke https://github.com/suzuki-shunsuke"
	app.Usage = "binary version control system"
	app.Commands = []cli.Command{
		handler.InitCommand,
		handler.InstallCommand,
	}
	app.Run(os.Args)
}
