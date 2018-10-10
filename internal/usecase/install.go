package usecase

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
	"github.com/suzuki-shunsuke/akoi/internal/util"
)

const (
	keyWordAnsible = "ansible"
)

// Install intalls binraries.
func Install(
	ctx context.Context, params domain.InstallParams,
	methods domain.InstallMethods,
) domain.Result {
	// suppress output log by third party library
	// https://github.com/joeybloggs/go-download/blob/b655936947da12d76bee4fa3b6af41a98db23e6f/download.go#L119
	log.SetOutput(methods.NewLoggerOutput())

	result := domain.Result{
		Packages: map[string]domain.PackageResult{}}
	if err := util.ValidateStruct(methods); err != nil {
		if methods.Fprintln != nil {
			methods.Fprintln(os.Stderr, err)
		}
		result.Msg = err.Error()
		return result
	}
	cfg, err := methods.ReadConfigFile(params.ConfigFilePath)
	if err != nil {
		methods.Fprintln(os.Stderr, err)
		result.Msg = err.Error()
		return result
	}
	cfg, err = setupConfig(cfg, methods)
	if err != nil {
		methods.Fprintln(os.Stderr, err)
		result.Msg = err.Error()
		return result
	}
	numOfPkgs := len(cfg.Packages)
	if numOfPkgs == 0 {
		return result
	}
	var wg sync.WaitGroup
	pkgResultChan := make(chan domain.PackageResult, numOfPkgs)
	for _, pkg := range cfg.Packages {
		wg.Add(1)
		go func(pkg domain.Package) {
			defer wg.Done()
			c, cancel := context.WithCancel(ctx)
			defer cancel()
			pkg = installPackage(c, pkg, params, methods)
			pkgResult := pkg.Result
			if pkgResult == nil {
				pkgResult = &domain.PackageResult{Name: pkg.Name}
			}
			for _, file := range pkg.Files {
				fileResult := file.Result
				pkgResult.Files[file.Name] = *fileResult
			}
			pkgResultChan <- *pkgResult
		}(pkg)
	}
	wg.Wait()
	close(pkgResultChan)
	for pkgResult := range pkgResultChan {
		result.Packages[pkgResult.Name] = pkgResult
	}
	return result
}
