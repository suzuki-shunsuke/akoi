package usecase

import (
	"github.com/suzuki-shunsuke/akoi/domain"
)

func installPackage(pkg *domain.Package, params *domain.InstallParams, methods *domain.InstallMethods) *domain.PackageResult {
	pkgResult := &domain.PackageResult{
		Files:   []domain.FileResult{},
		Version: pkg.Version,
		Name:    pkg.Name,
	}
	for _, file := range pkg.Files {
		fileResult, err := installFile(pkg, &file, params, methods)
		if fileResult == nil {
			fileResult = &domain.FileResult{}
		}
		if fileResult.Changed {
			pkgResult.Changed = true
		}
		if err != nil {
			if fileResult.Error == "" {
				fileResult.Error = err.Error()
			}
			pkgResult.Files = append(pkgResult.Files, *fileResult)
			continue
		}
		fr, err := createLink(pkg, &file, params, methods)
		if fr == nil {
			fr = &domain.FileResult{}
		}
		if fr.Changed {
			pkgResult.Changed = true
		}
		if err != nil {
			if fileResult.Error == "" {
				fileResult.Error = err.Error()
			}
			pkgResult.Files = append(pkgResult.Files, *fileResult)
		}
		pkgResult.Files = append(pkgResult.Files, *fileResult)
	}
	return pkgResult
}
