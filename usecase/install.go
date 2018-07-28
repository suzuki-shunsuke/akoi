package usecase

import (
	"fmt"
	"net/url"
	"path/filepath"
	"text/template"

	"github.com/suzuki-shunsuke/akoi/domain"
	"github.com/suzuki-shunsuke/akoi/util"
)

const (
	keyWordAnsible = "ansible"
)

// setupConfig compiles and renders templates of domain.Config .
func setupConfig(cfg *domain.Config, methods *domain.InstallMethods) error {
	tpl, err := template.New("bin_path").Parse(cfg.BinPath)
	if err != nil {
		return err
	}
	cfg.BinPathTpl = tpl

	tpl, err = template.New("link_path").Parse(cfg.LinkPath)
	if err != nil {
		return err
	}
	cfg.LinkPathTpl = tpl

	for pkgName, pkg := range cfg.Packages {
		pkg.Name = pkgName
		tpl, err := template.New("pkg_url").Parse(pkg.RawURL)
		if err != nil {
			return err
		}
		u, err := util.RenderTpl(tpl, pkg)
		if err != nil {
			return err
		}
		u2, err := url.Parse(u)
		if err != nil {
			return err
		}
		pkg.URL = u2
		pkg.Archiver = methods.GetArchiver(u2.Path)
		for i, file := range pkg.Files {
			dst, err := util.RenderTpl(
				cfg.BinPathTpl, &domain.TemplateParams{
					Name: file.Name, Version: pkg.Version,
				})
			if err != nil {
				return err
			}
			file.Bin = dst

			lnPath, err := util.RenderTpl(
				cfg.LinkPathTpl, &domain.TemplateParams{
					Name: file.Name, Version: pkg.Version,
				})
			if err != nil {
				return err
			}
			file.Link = lnPath
			pkg.Files[i] = file
		}
		cfg.Packages[pkgName] = pkg
	}

	return nil
}

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
	if err := setupConfig(cfg, methods); err != nil {
		return result, err
	}
	for _, pkg := range cfg.Packages {
		// TODO goroutine
		pkgResult, err := installPackage(&pkg, params, methods)
		if pkgResult == nil {
			pkgResult = &domain.PackageResult{}
		}
		result.Packages[pkg.Name] = *pkgResult
		if pkgResult.Changed {
			result.Changed = true
		}
		if err != nil {
			if pkgResult.Error == "" {
				pkgResult.Error = err.Error()
			}
			result.Packages[pkg.Name] = *pkgResult
			return result, err
		}
	}
	return result, nil
}

func createLink(dst string, pkg *domain.Package, file *domain.File, params *domain.InstallParams, methods *domain.InstallMethods) (*domain.FileResult, error) {
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

func installFile(dst string, pkg *domain.Package, file *domain.File, params *domain.InstallParams, methods *domain.InstallMethods) (*domain.FileResult, error) {
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

	ustr := pkg.URL.String()
	if params.Format != keyWordAnsible {
		fmt.Printf("downloading %s: %s\n", pkg.Name, ustr)
	}
	resp, err := methods.Download(ustr)
	if err != nil {
		return fileResult, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fileResult, fmt.Errorf(
			"failed to download %s from %s: %d", pkg.Name, ustr, resp.StatusCode)
	}
	tmpDir, err := methods.TempDir()
	if err != nil {
		return fileResult, err
	}
	defer methods.RemoveAll(tmpDir)
	arc := pkg.Archiver
	// TODO support not archived file
	if arc != nil {
		if params.Format != keyWordAnsible {
			fmt.Printf("unarchive %s\n", pkg.Name)
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

func installPackage(pkg *domain.Package, params *domain.InstallParams, methods *domain.InstallMethods) (*domain.PackageResult, error) {
	pkgResult := &domain.PackageResult{
		Files:   []domain.FileResult{},
		Version: pkg.Version,
	}
	for _, file := range pkg.Files {
		fileResult, err := installFile(file.Bin, pkg, &file, params, methods)
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
		fr, err := createLink(file.Bin, pkg, &file, params, methods)
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
