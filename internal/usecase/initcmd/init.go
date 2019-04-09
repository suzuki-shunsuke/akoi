package initcmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
)

// InitConfigFile creates a configuration file if it doesn't exist.
func InitConfigFile(params *domain.InitParams, fsys domain.FileSystem) error {
	dest := params.Dest
	if fsys.ExistFile(dest) {
		return nil
	}
	dir := filepath.Dir(dest)
	if !fsys.ExistFile(dir) {
		fmt.Printf("create a directory %s\n", dir)
		if err := fsys.MkdirAll(dir); err != nil {
			return err
		}
	}
	fmt.Printf("create %s\n", dest)
	return fsys.WriteFile(dest, []byte(strings.Trim(domain.ConfigTemplate, "\n")))
}
