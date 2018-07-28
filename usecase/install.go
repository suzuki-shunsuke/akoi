package usecase

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sync"
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
		pkg.Archiver = methods.GetArchiver(u2.Path, pkg.ArchiveType)
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
			pkgResult := installPackage(&pkg, params, methods)
			if pkgResult == nil {
				pkgResult = &domain.PackageResult{Name: pkg.Name}
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
		if pkgResult.Error != "" {
			result.Failed = true
		}
	}
	return result
}

func createLink(pkg *domain.Package, file *domain.File, params *domain.InstallParams, methods *domain.InstallMethods) (*domain.FileResult, error) {
	fileResult := &domain.FileResult{}
	linkRelPath, err := filepath.Rel(filepath.Dir(file.Link), file.Bin)
	if err != nil {
		if params.Format != keyWordAnsible {
			fmt.Fprintln(os.Stderr, err)
		}
		return fileResult, err
	}
	if fi, err := methods.GetFileLstat(file.Link); err == nil {
		switch mode := fi.Mode(); {
		case mode.IsDir():
			if params.Format != keyWordAnsible {
				fmt.Fprintf(os.Stderr, "%s has already existed and is a directory\n", file.Link)
			}
			return fileResult, fmt.Errorf("%s has already existed and is a directory", file.Link)
		case mode&os.ModeNamedPipe != 0:
			if params.Format != keyWordAnsible {
				fmt.Fprintf(os.Stderr, "%s has already existed and is a named pipe\n", file.Link)
			}
			return fileResult, fmt.Errorf("%s has already existed and is a named pipe", file.Link)
		case mode.IsRegular():
			if params.Format != keyWordAnsible {
				fmt.Printf("remove %s\n", file.Link)
			}
			if err := methods.RemoveFile(file.Link); err != nil {
				if params.Format != keyWordAnsible {
					fmt.Fprintln(os.Stderr, err)
				}
				return fileResult, err
			}
			fileResult.Changed = true
			if params.Format != keyWordAnsible {
				fmt.Printf("create link %s -> %s\n", file.Link, linkRelPath)
			}
			if err := methods.MkLink(linkRelPath, file.Link); err != nil {
				if params.Format != keyWordAnsible {
					fmt.Fprintln(os.Stderr, err)
				}
				return fileResult, err
			}
			fileResult.Changed = true
			return fileResult, nil
		case mode&os.ModeSymlink != 0:
			lnDest, err := methods.ReadLink(file.Link)
			if err != nil {
				if params.Format != keyWordAnsible {
					fmt.Fprintln(os.Stderr, err)
				}
				return fileResult, err
			}
			if linkRelPath == lnDest {
				return fileResult, nil
			}
			if params.Format != keyWordAnsible {
				fmt.Printf("remove link %s -> %s\n", file.Link, lnDest)
			}
			if err := methods.RemoveLink(file.Link); err != nil {
				if params.Format != keyWordAnsible {
					fmt.Fprintln(os.Stderr, err)
				}
				return fileResult, err
			}
			fileResult.Changed = true
			if params.Format != keyWordAnsible {
				fmt.Printf("create link %s -> %s\n", file.Link, linkRelPath)
			}
			if err := methods.MkLink(linkRelPath, file.Link); err != nil {
				if params.Format != keyWordAnsible {
					fmt.Fprintln(os.Stderr, err)
				}
				return fileResult, err
			}
			return fileResult, nil
		default:
			return fileResult, fmt.Errorf("unexpected file mode %s: %s", file.Link, mode.String())
		}
	}
	if params.Format != keyWordAnsible {
		fmt.Printf("create link %s -> %s\n", file.Link, linkRelPath)
	}
	if err := methods.MkLink(linkRelPath, file.Link); err != nil {
		if params.Format != keyWordAnsible {
			fmt.Fprintln(os.Stderr, err)
		}
		return fileResult, err
	}
	fileResult.Changed = true
	return fileResult, nil
}

func installFile(pkg *domain.Package, file *domain.File, params *domain.InstallParams, methods *domain.InstallMethods) (*domain.FileResult, error) {
	dst := file.Bin
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
		err := methods.Chmod(dst, mode)
		if err != nil && params.Format != keyWordAnsible {
			fmt.Fprintln(os.Stderr, err)
		}
		return fileResult, err
	}

	ustr := pkg.URL.String()
	if params.Format != keyWordAnsible {
		fmt.Printf("downloading %s: %s\n", pkg.Name, ustr)
	}
	resp, err := methods.Download(ustr)
	if err != nil {
		if params.Format != keyWordAnsible {
			fmt.Fprintln(os.Stderr, err)
		}
		return fileResult, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		if params.Format != keyWordAnsible {
			fmt.Fprintf(os.Stderr, "failed to download %s from %s: %d\n", pkg.Name, ustr, resp.StatusCode)
		}
		return fileResult, fmt.Errorf(
			"failed to download %s from %s: %d", pkg.Name, ustr, resp.StatusCode)
	}
	tmpDir, err := methods.TempDir()
	if err != nil {
		if params.Format != keyWordAnsible {
			fmt.Fprintln(os.Stderr, err)
		}
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
			if params.Format != keyWordAnsible {
				fmt.Fprintln(os.Stderr, err)
			}
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
						if params.Format != keyWordAnsible {
							fmt.Fprintln(os.Stderr, err)
						}
						return fileResult, err
					}
					fileResult.Changed = true
				}
				if params.Format != keyWordAnsible {
					fmt.Printf("install %s\n", dst)
				}
				if err := methods.CopyFile(filepath.Join(tmpDir, f.Archive), dst); err != nil {
					if params.Format != keyWordAnsible {
						fmt.Fprintln(os.Stderr, err)
					}
					return fileResult, err
				}
				fileResult.Changed = true
			}
			if err == nil && fi.Mode() == mode {
				continue
			}
			if err := methods.Chmod(dst, mode); err != nil {
				if params.Format != keyWordAnsible {
					fmt.Fprintln(os.Stderr, err)
				}
				return fileResult, err
			}
			fileResult.Changed = true
		}
	}
	return fileResult, nil
}

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
