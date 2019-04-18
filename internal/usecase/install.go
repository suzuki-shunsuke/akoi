package usecase

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/suzuki-shunsuke/gomic/gomic"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
)

const (
	keyWordAnsible = "ansible"
)

func (lgc *Logic) Install(
	ctx context.Context, params domain.InstallParams,
	cfgReader domain.ConfigReader, getArchiver domain.GetArchiver,
	downloader domain.Downloader, getGzipReader domain.GetGzipReader,
) domain.Result {
	// suppress output log by third party library
	// https://github.com/joeybloggs/go-download/blob/b655936947da12d76bee4fa3b6af41a98db23e6f/download.go#L119
	log.SetOutput(NewWriter(nil, gomic.DoNothing))

	result := domain.Result{
		Packages: map[string]domain.PackageResult{}}
	cfg, err := cfgReader.Read(params.ConfigFilePath)
	if err != nil {
		lgc.Printer.Fprintln(os.Stderr, err)
		result.Msg = err.Error()
		return result
	}
	cfg, err = lgc.Logic.SetupConfig(cfg, getArchiver)
	if err != nil {
		lgc.Printer.Fprintln(os.Stderr, err)
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
			pkg = lgc.Logic.InstallPackage(ctx, pkg, params, downloader, getGzipReader)
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
