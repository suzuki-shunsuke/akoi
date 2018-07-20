package command

import (
	"github.com/urfave/cli"

	"github.com/suzuki-shunsuke/akoi/domain"
	"github.com/suzuki-shunsuke/akoi/registry"
	"github.com/suzuki-shunsuke/akoi/usecase"
)

// Init is the sub command "init".
func Init(c *cli.Context) error {
	err := usecase.InitConfigFile(
		&domain.InitParams{
			Dest: c.String("dest"),
		},
		registry.NewInitMethods(),
	)
	if err == nil {
		return nil
	}
	return cli.NewExitError(err.Error(), 1)
}
