package handler

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
	"github.com/suzuki-shunsuke/akoi/internal/infra"
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

	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan, syscall.SIGHUP, syscall.SIGINT,
		syscall.SIGTERM, syscall.SIGQUIT)
	resultChan := make(chan domain.Result)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logic := newLogic()
	go func() {
		resultChan <- logic.Install(
			ctx, params, infra.Printer{},
			infra.ConfigReader{}, infra.GetArchiver{}, infra.Downloader{}, infra.GetGzipReader{})
	}()
	select {
	case result := <-resultChan:
		close(signalChan)
		if !result.Failed() {
			s := result.String(params.Format)
			if s != "" {
				fmt.Println(s)
			}
			return nil
		}
		return cli.NewExitError(result.String(params.Format), 1)
	case sig := <-signalChan:
		close(resultChan)
		return cli.NewExitError(sig.String(), 1)
	}
}

func newLogic() domain.Logic {
	lgc := &usecase.Logic{
		Fsys: infra.FileSystem{},
	}
	lgc.Logic = lgc
	return lgc
}
