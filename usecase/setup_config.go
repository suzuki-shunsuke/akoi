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
func setupConfig(cfg *domain.Config, methods *domain.SetupConfigMethods) error {
	cfgBinDirTpl, err := template.New("cfg_bin_dir").Parse(cfg.BinDirTplStr)
	if err != nil {
		return err
	}

	cfgLinkDirTpl, err := template.New("cfg_link_dir").Parse(cfg.LinkDirTplStr)
	if err != nil {
		return err
	}

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

		if pkg.LinkDirTplStr == "" {
			pkg.LinkDirTpl = cfgLinkDirTpl
		} else {
			tpl, err := template.New("pkg_link_dir").Parse(pkg.LinkDirTplStr)
			if err != nil {
				return err
			}
			pkg.LinkDirTpl = tpl
		}

		if pkg.BinDirTplStr == "" {
			pkg.BinDirTpl = cfgBinDirTpl
		} else {
			tpl, err := template.New("pkg_bin_dir").Parse(pkg.BinDirTplStr)
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
		// TODO validate archiver
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

			if file.LinkDirTplStr == "" {
				file.LinkDirTpl = pkg.LinkDirTpl
			} else {
				tpl, err := template.New("file_link_dir").Parse(file.LinkDirTplStr)
				if err != nil {
					return err
				}
				file.LinkDirTpl = tpl
			}

			if file.BinDirTplStr == "" {
				file.BinDirTpl = pkg.BinDirTpl
			} else {
				tpl, err := template.New("file_bin_dir").Parse(file.BinDirTplStr)
				if err != nil {
					return err
				}
				file.BinDirTpl = tpl
			}

			if file.BinSeparator == "" {
				file.BinSeparator = pkg.BinSeparator
			}

			file.BinDir, err = util.RenderTpl(
				file.BinDirTpl, &domain.TemplateParams{
					Name: file.Name, Version: pkg.Version,
				})
			if err != nil {
				return err
			}
			file.Bin = filepath.Join(
				file.BinDir, fmt.Sprintf("%s%s%s", file.Name, file.BinSeparator, pkg.Version))

			lnPath, err := util.RenderTpl(
				file.LinkDirTpl, &domain.TemplateParams{
					Name: file.Name, Version: pkg.Version,
				})
			if err != nil {
				return err
			}
			file.Link = filepath.Join(lnPath, file.Name)
			pkg.Files[i] = file
		}
		cfg.Packages[pkgName] = pkg
	}

	return nil
}
