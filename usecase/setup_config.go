package usecase

import (
	"fmt"
	"net/url"
	"path/filepath"
	"text/template"

	"github.com/suzuki-shunsuke/akoi/domain"
	"github.com/suzuki-shunsuke/akoi/util"
)

// setupConfig compiles and renders templates of domain.Config .
func setupConfig(cfg *domain.Config, methods *domain.InstallMethods) error {
	tpl, err := template.New("cfg_bin_dir").Parse(cfg.BinDir)
	if err != nil {
		return err
	}
	cfg.BinDirTpl = tpl

	tpl, err = template.New("cfg_link_dir").Parse(cfg.LinkDir)
	if err != nil {
		return err
	}
	cfg.LinkDirTpl = tpl

	if cfg.BinSeparator == "" {
		cfg.BinSeparator = "-"
	}

	for pkgName, pkg := range cfg.Packages {
		if pkg.Result == nil {
			pkg.Result = &domain.PackageResult{
				Name:  pkgName,
				Files: map[string]domain.FileResult{},
			}
		}

		if pkg.LinkDir == "" {
			pkg.LinkDir = cfg.LinkDir
			pkg.LinkDirTpl = cfg.LinkDirTpl
		} else {
			tpl, err := template.New("pkg_link_dir").Parse(pkg.LinkDir)
			if err != nil {
				return err
			}
			pkg.LinkDirTpl = tpl
		}

		if pkg.BinDir == "" {
			pkg.BinDir = cfg.BinDir
			pkg.BinDirTpl = cfg.BinDirTpl
		} else {
			tpl, err := template.New("pkg_bin_dir").Parse(pkg.BinDir)
			if err != nil {
				return err
			}
			pkg.BinDirTpl = tpl
		}

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
		if pkg.BinSeparator == "" {
			pkg.BinSeparator = cfg.BinSeparator
		}
		for i, file := range pkg.Files {
			if file.Result == nil {
				file.Result = &domain.FileResult{}
			}
			if file.Mode == 0 {
				file.Mode = 0755
			}

			if file.LinkDir == "" {
				file.LinkDir = pkg.LinkDir
				file.LinkDirTpl = pkg.LinkDirTpl
			} else {
				tpl, err := template.New("file_link_dir").Parse(file.LinkDir)
				if err != nil {
					return err
				}
				file.LinkDirTpl = tpl
			}

			if file.BinDir == "" {
				file.BinDir = pkg.BinDir
				file.BinDirTpl = pkg.BinDirTpl
			} else {
				tpl, err := template.New("file_bin_dir").Parse(file.BinDir)
				if err != nil {
					return err
				}
				file.BinDirTpl = tpl
			}

			if file.BinSeparator == "" {
				file.BinSeparator = pkg.BinSeparator
			}

			dst, err := util.RenderTpl(
				file.BinDirTpl, &domain.TemplateParams{
					Name: file.Name, Version: pkg.Version,
				})
			if err != nil {
				return err
			}
			file.Bin = filepath.Join(
				dst, fmt.Sprintf("%s%s%s", file.Name, file.BinSeparator, pkg.Version))

			lnPath, err := util.RenderTpl(
				file.LinkDirTpl, &domain.TemplateParams{
					Name: file.Name, Version: pkg.Version,
				})
			if err != nil {
				return err
			}
			file.Link = fmt.Sprintf("%s%s", lnPath, file.Name)
			file.Link = filepath.Join(lnPath, file.Name)
			pkg.Files[i] = file
		}
		cfg.Packages[pkgName] = pkg
	}

	return nil
}
