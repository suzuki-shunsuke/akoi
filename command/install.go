package command

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/suzuki-shunsuke/akoi/domain"
	"github.com/suzuki-shunsuke/akoi/registry"
	"github.com/suzuki-shunsuke/akoi/usecase"
)

// Install is the sub command "install".
func Install(c *cli.Context) error {
	params := &domain.InstallParams{
		ConfigFilePath: c.String("config"),
		Format:         c.String("format"),
	}
	result, err := usecase.Install(params, registry.NewInstallMethods())
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
