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
) (domain.Result, error) {
	// suppress output log by third party library
	// https://github.com/joeybloggs/go-download/blob/b655936947da12d76bee4fa3b6af41a98db23e6f/download.go#L119
	log.SetOutput(NewWriter(nil, gomic.DoNothing))

	result := domain.Result{
		Packages: map[string]domain.PackageResult{}}
	cfg, err := lgc.CfgReader.Read(params.ConfigFilePath)
	if err != nil {
		lgc.Printer.Fprintln(os.Stderr, err)
		result.Msg = err.Error()
		return result, err
	}
	cfg, err = lgc.Logic.SetupConfig(cfg, params.ConfigFilePath)
	if err != nil {
		lgc.Printer.Fprintln(os.Stderr, err)
		result.Msg = err.Error()
		return result, err
	}
	numOfPkgs := len(cfg.Packages)
	if numOfPkgs == 0 {
		return result, nil
	}
	var wg sync.WaitGroup
	pkgResultChan := make(chan domain.PackageResult, numOfPkgs)
	if cfg.MaxParallelDownloadCount != 0 {
		lgc.maxParallelDownloadCountChan = make(
			chan struct{}, cfg.MaxParallelDownloadCount)
	}
	for _, pkg := range cfg.Packages {
		wg.Add(1)
		go func(pkg domain.Package) {
			defer wg.Done()
			p, err := lgc.Logic.InstallPackage(ctx, pkg, params)
			if err != nil {
				p.Result.Error = err.Error()
			}
			pkgResult := p.Result
			if pkgResult == nil {
				pkgResult = &domain.PackageResult{Name: p.Name}
			}
			for _, file := range p.Files {
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
	return result, nil
}

func (lgc *Logic) pushMaxParallelDownloadCount() {
	if lgc.maxParallelDownloadCountChan != nil {
		lgc.maxParallelDownloadCountChan <- struct{}{}
	}
}

func (lgc *Logic) popMaxParallelDownloadCount() {
	if lgc.maxParallelDownloadCountChan != nil {
		<-lgc.maxParallelDownloadCountChan
	}
}
