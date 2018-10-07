package initcmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
	"github.com/suzuki-shunsuke/akoi/internal/util"
)

// InitConfigFile creates a configuration file if it doesn't exist.
func InitConfigFile(params *domain.InitParams, methods *domain.InitMethods) error {
	if err := util.ValidateStruct(methods); err != nil {
		return err
	}
	dest := params.Dest
	if methods.Exist(dest) {
		return nil
	}
	dir := filepath.Dir(dest)
	if !methods.Exist(dir) {
		fmt.Printf("create a directory %s\n", dir)
		if err := methods.MkdirAll(dir); err != nil {
			return err
		}
	}
	fmt.Printf("create %s\n", dest)
	return methods.Write(dest, []byte(strings.Trim(domain.ConfigTemplate, "\n")))
}
