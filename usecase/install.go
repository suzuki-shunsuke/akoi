package usecase

import (
	"sync"

	"github.com/suzuki-shunsuke/akoi/domain"
	"github.com/suzuki-shunsuke/akoi/util"
)

const (
	keyWordAnsible = "ansible"
)

// Install intalls binraries.
func Install(params *domain.InstallParams, methods *domain.InstallMethods) *domain.Result {
	result := &domain.Result{
		Packages: map[string]domain.PackageResult{}}
	if err := util.ValidateStruct(methods); err != nil {
		result.Msg = err.Error()
		result.Failed = true
		return result
	}
	cfg, err := methods.ReadConfigFile(params.ConfigFilePath)
	if err != nil {
		result.Msg = err.Error()
		result.Failed = true
		return result
	}
	if err := setupConfig(cfg, methods); err != nil {
		result.Msg = err.Error()
		result.Failed = true
		return result
	}
	numOfPkgs := len(cfg.Packages)
	if numOfPkgs == 0 {
		return result
	}
	var wg sync.WaitGroup
	pkgResultChan := make(chan domain.PackageResult, numOfPkgs)
	for _, pkg := range cfg.Packages {
		// TODO goroutine
		wg.Add(1)
		go func(pkg domain.Package) {
			defer wg.Done()
			installPackage(&pkg, params, methods)
			pkgResult := pkg.Result
			if pkgResult == nil {
				pkgResult = &domain.PackageResult{Name: pkg.Name}
			}
			for _, file := range pkg.Files {
				fileResult := file.Result
				if fileResult.Changed {
					pkgResult.Changed = true
				}
				if fileResult.Error != "" {
					pkgResult.Failed = true
				}
				pkgResult.Files[file.Name] = *fileResult
			}
			pkgResultChan <- *pkgResult
		}(pkg)
	}
	wg.Wait()
	close(pkgResultChan)
	for pkgResult := range pkgResultChan {
		result.Packages[pkgResult.Name] = pkgResult
		if pkgResult.Changed {
			result.Changed = true
		}
		if pkgResult.Failed {
			result.Failed = true
		}
	}
	return result
}
