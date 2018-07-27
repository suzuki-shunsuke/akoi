package usecase

import (
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/suzuki-shunsuke/akoi/domain"
	"github.com/suzuki-shunsuke/akoi/util"
)

const (
	keyWordAnsible = "ansible"
)

// Install intalls binraries.
func Install(params *domain.InstallParams, methods *domain.InstallMethods) (*domain.Result, error) {
	result := &domain.Result{
		Packages: map[string]domain.PackageResult{}}
	if err := util.ValidateStruct(methods); err != nil {
		return result, err
	}
	cfg, err := methods.ReadConfigFile(params.ConfigFilePath)
	if err != nil {
		return result, err
	}
	if err := cfg.Setup(); err != nil {
		return result, err
	}
	for pkgName, pkg := range cfg.Packages {
		// TODO goroutine
		pkgResult, err := installPackage(pkgName, &pkg, params, methods)
		if pkgResult == nil {
			pkgResult = &domain.PackageResult{}
		}
		result.Packages[pkgName] = *pkgResult
		if pkgResult.Changed {
			result.Changed = true
		}
		if err != nil {
			if pkgResult.Error == "" {
				pkgResult.Error = err.Error()
			}
			result.Packages[pkgName] = *pkgResult
			return result, err
		}
	}
	return result, nil
}

func createLink(pkgName, dst string, pkg *domain.Package, file *domain.File, params *domain.InstallParams, methods *domain.InstallMethods) (*domain.FileResult, error) {
	fileResult := &domain.FileResult{}
	if _, err := methods.GetFileLstat(file.Link); err != nil {
		if _, err := methods.GetFileStat(file.Link); err == nil {
			// TODO force remove option
			return fileResult, fmt.Errorf("%s has already existed and is not a symbolic link", file.Link)
		}
		p, err := filepath.Rel(filepath.Dir(file.Link), dst)
		if err != nil {
			return fileResult, err
		}
		if params.Format != keyWordAnsible {
			fmt.Printf("create link %s -> %s\n", file.Link, p)
		}
		if err := methods.MkLink(p, file.Link); err != nil {
			return fileResult, err
		}
		fileResult.Changed = true
		return fileResult, nil
	}
	lnDest, err := methods.ReadLink(file.Link)
	if err != nil {
		return fileResult, err
	}
	p, err := filepath.Rel(filepath.Dir(file.Link), dst)
	if err != nil {
		return fileResult, err
	}
	if p == lnDest {
		return fileResult, nil
	}
	if params.Format != keyWordAnsible {
		fmt.Printf("remove link %s -> %s\n", file.Link, lnDest)
	}
	if err := methods.RemoveLink(file.Link); err != nil {
		return fileResult, err
	}
	fileResult.Changed = true
	if params.Format != keyWordAnsible {
		fmt.Printf("create link %s -> %s\n", file.Link, p)
	}
	if err := methods.MkLink(p, file.Link); err != nil {
		return fileResult, err
	}
	return fileResult, nil
}

func installFile(pkgName, dst string, pkg *domain.Package, file *domain.File, params *domain.InstallParams, methods *domain.InstallMethods) (*domain.FileResult, error) {
	fileResult := &domain.FileResult{
		Name: file.Name,
	}
	mode := file.Mode
	if mode == 0 {
		mode = 0755
	}
	if fi, err := methods.GetFileStat(dst); err == nil {
		if fi.Mode() == mode {
			return fileResult, nil
		}
		if params.Format != keyWordAnsible {
			fmt.Printf("chmod %s %s\n", mode.String(), dst)
		}
		return fileResult, methods.Chmod(dst, mode)
	}

	u := pkg.GetURL()
	if params.Format != keyWordAnsible {
		fmt.Printf("downloading %s: %s\n", pkgName, u)
	}
	resp, err := methods.Download(u)
	if err != nil {
		return fileResult, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fileResult, fmt.Errorf("failed to download %s from %s: %d", pkgName, u, resp.StatusCode)
	}
	tmpDir, err := methods.TempDir()
	if err != nil {
		return fileResult, err
	}
	defer methods.RemoveAll(tmpDir)
	u2, err := url.Parse(u)
	if err != nil {
		return fileResult, err
	}
	arc := methods.GetArchiver(u2.Path)
	// TODO support not archived file
	if arc != nil {
		if params.Format != keyWordAnsible {
			fmt.Printf("unarchive %s\n", pkgName)
		}
		if err := arc.Read(resp.Body, tmpDir); err != nil {
			return fileResult, err
		}
		for _, f := range pkg.Files {
			fi, err := methods.GetFileStat(f.Bin)
			if err != nil {
				dir := filepath.Dir(f.Bin)
				if _, err := methods.GetFileStat(dir); err != nil {
					if params.Format != keyWordAnsible {
						fmt.Printf("create directory %s\n", dir)
					}
					if err := methods.MkdirAll(dir); err != nil {
						return fileResult, err
					}
					fileResult.Changed = true
				}
				if params.Format != keyWordAnsible {
					fmt.Printf("install %s\n", dst)
				}
				if err := methods.CopyFile(filepath.Join(tmpDir, f.Archive), dst); err != nil {
					return fileResult, err
				}
				fileResult.Changed = true
			}
			if err == nil && fi.Mode() == mode {
				continue
			}
			if err := methods.Chmod(dst, mode); err != nil {
				return fileResult, err
			}
			fileResult.Changed = true
		}
	}
	return fileResult, nil
}

func installPackage(pkgName string, pkg *domain.Package, params *domain.InstallParams, methods *domain.InstallMethods) (*domain.PackageResult, error) {
	pkgResult := &domain.PackageResult{
		Files:   []domain.FileResult{},
		Version: pkg.Version,
	}
	for _, file := range pkg.Files {
		fileResult, err := installFile(
			pkgName, file.Bin, pkg, &file, params, methods)
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
			return pkgResult, err
		}
		fr, err := createLink(
			pkgName, file.Bin, pkg, &file, params, methods)
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
			return pkgResult, err
		}
		pkgResult.Files = append(pkgResult.Files, *fileResult)
		continue
	}
	return pkgResult, nil
}
