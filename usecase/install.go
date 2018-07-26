package usecase

import (
	"fmt"
	"path/filepath"
	"text/template"

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
	binPathTpl, err := cfg.GetBinPathTpl()
	if err != nil {
		return result, err
	}
	linkPathTpl, err := cfg.GetLinkPathTpl()
	if err != nil {
		return result, err
	}
	for pkgName, pkg := range cfg.Packages {
		// TODO goroutine
		pkgResult, err := installPackage(
			pkgName, &pkg, params, methods, binPathTpl, linkPathTpl)
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

func createLink(pkgName, dst string, pkg *domain.Package, file *domain.File, params *domain.InstallParams, methods *domain.InstallMethods, linkPathTpl *template.Template) (*domain.FileResult, error) {
	fileResult := &domain.FileResult{}
	lnPath, err := util.RenderTpl(
		linkPathTpl, &domain.TemplateParams{
			Name: file.Name, Version: pkg.Version,
		})
	if err != nil {
		return fileResult, err
	}

	if _, err := methods.GetFileLstat(lnPath); err != nil {
		if _, err := methods.GetFileStat(lnPath); err == nil {
			// TODO force remove option
			return fileResult, fmt.Errorf("%s has already existed and is not a symbolic link", lnPath)
		}
		p, err := filepath.Rel(filepath.Dir(lnPath), dst)
		if err != nil {
			return fileResult, err
		}
		if params.Format != keyWordAnsible {
			fmt.Printf("create link %s -> %s\n", lnPath, p)
		}
		if err := methods.MkLink(p, lnPath); err != nil {
			return fileResult, err
		}
		fileResult.Changed = true
		return fileResult, nil
	}
	lnDest, err := methods.ReadLink(lnPath)
	if err != nil {
		return fileResult, err
	}
	p, err := filepath.Rel(filepath.Dir(lnPath), dst)
	if err != nil {
		return fileResult, err
	}
	if p == lnDest {
		return fileResult, nil
	}
	if params.Format != keyWordAnsible {
		fmt.Printf("remove link %s -> %s\n", lnPath, lnDest)
	}
	if err := methods.RemoveLink(lnPath); err != nil {
		return fileResult, err
	}
	fileResult.Changed = true
	if params.Format != keyWordAnsible {
		fmt.Printf("create link %s -> %s\n", lnPath, p)
	}
	if err := methods.MkLink(p, lnPath); err != nil {
		return fileResult, err
	}
	return fileResult, nil
}

func installFile(pkgName, dst string, pkg *domain.Package, file *domain.File, params *domain.InstallParams, methods *domain.InstallMethods, binPathTpl *template.Template) (*domain.FileResult, error) {
	fileResult := &domain.FileResult{
		Name: file.Name,
	}
	if fi, err := methods.GetFileStat(dst); err == nil {
		if fi.Mode() == 0755 {
			return fileResult, nil
		}
		return fileResult, methods.Chmod(dst, 0755)
	}

	u, err := pkg.GetURL()
	if err != nil {
		return fileResult, err
	}
	if params.Format != keyWordAnsible {
		fmt.Printf("downloading %s: %s\n", pkgName, u)
	}
	resp, err := methods.Download(u)
	if err != nil {
		return fileResult, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fileResult, fmt.Errorf("failed to download %s from %s: %d", pkgName, pkg.URL, resp.StatusCode)
	}
	tmpDir, err := methods.TempDir()
	if err != nil {
		return fileResult, err
	}
	defer methods.RemoveAll(tmpDir)
	arc, err := pkg.GetArchiver()
	if err != nil {
		return fileResult, err
	}
	// TODO support not archived file
	if arc != nil {
		if params.Format != keyWordAnsible {
			fmt.Printf("unarchive %s\n", pkgName)
		}
		if err := arc.Read(resp.Body, tmpDir); err != nil {
			return fileResult, err
		}
		for _, f := range pkg.Files {
			dst, err := util.RenderTpl(
				binPathTpl, &domain.TemplateParams{
					Name: f.Name, Version: pkg.Version,
				})
			if err != nil {
				return fileResult, err
			}
			fi, err := methods.GetFileStat(dst)
			if err != nil {
				dir := filepath.Dir(dst)
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
			if err == nil && fi.Mode() == 0755 {
				continue
			}
			if err := methods.Chmod(dst, 0755); err != nil {
				return fileResult, err
			}
			fileResult.Changed = true
		}
	}
	return fileResult, nil
}

func installPackage(pkgName string, pkg *domain.Package, params *domain.InstallParams, methods *domain.InstallMethods, binPathTpl, linkPathTpl *template.Template) (*domain.PackageResult, error) {
	pkgResult := &domain.PackageResult{
		Files:   []domain.FileResult{},
		Version: pkg.Version,
	}
	for _, file := range pkg.Files {
		dst, err := util.RenderTpl(
			binPathTpl, &domain.TemplateParams{
				Name: file.Name, Version: pkg.Version,
			})
		if err != nil {
			return pkgResult, err
		}
		fileResult, err := installFile(
			pkgName, dst, pkg, &file, params, methods, binPathTpl)
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
			pkgName, dst, pkg, &file, params, methods, linkPathTpl)
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
