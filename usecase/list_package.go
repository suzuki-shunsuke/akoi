package usecase

import (
	"fmt"
	"path/filepath"

	"github.com/suzuki-shunsuke/akoi/domain"
)

func listPackage(pkg *domain.Package, params *domain.ListParams, methods *domain.ListMethods) {
	for _, file := range pkg.Files {
		fmt.Printf("%s/%s\n", pkg.Name, file.Name)
		s := filepath.Join(file.BinDir, fmt.Sprintf("%s%s", file.Name, file.BinSeparator))
		m, err := filepath.Glob(fmt.Sprintf("%s*", s))
		if err != nil {
			fmt.Println(err)
			return
		}
		n := len(s)
		for _, f := range m {
			fmt.Printf("  %s\n", f[n:])
		}
	}
}
