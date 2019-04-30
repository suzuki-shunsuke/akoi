package usecase

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/suzuki-shunsuke/go-cliutil"
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

	if params.ConfigFilePath == "" {
		wd, err := lgc.Fsys.Getwd()
		if err != nil {
			return result, err
		}
		params.ConfigFilePath, err = cliutil.FindFile(
			wd, ".akoi.yml", lgc.Fsys.ExistFile)
		if err != nil {
			params.ConfigFilePath = "/etc/akoi/akoi.yml"
			if !lgc.Fsys.ExistFile(params.ConfigFilePath) {
				return result, fmt.Errorf("configuration file is not found")
			}
		}
	} else {
		if !lgc.Fsys.ExistFile(params.ConfigFilePath) {
			return result, fmt.Errorf(
				"configuration file is not found: %s", params.ConfigFilePath)
		}
	}

	cfg, err := lgc.CfgReader.Read(params.ConfigFilePath)
	if err != nil {
		lgc.Printer.Fprintln(os.Stderr, err)
		return result, err
	}
	cfg, err = lgc.Logic.SetupConfig(cfg, params.ConfigFilePath)
	if err != nil {
		lgc.Printer.Fprintln(os.Stderr, err)
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
