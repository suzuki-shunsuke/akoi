package usecase

import (
	"net/url"
	"text/template"

	"github.com/suzuki-shunsuke/akoi/domain"
	"github.com/suzuki-shunsuke/akoi/util"
)

// setupConfig compiles and renders templates of domain.Config .
func setupConfig(cfg *domain.Config, methods *domain.InstallMethods) error {
	tpl, err := template.New("cfg_bin_path").Parse(cfg.BinPath)
	if err != nil {
		return err
	}
	cfg.BinPathTpl = tpl

	tpl, err = template.New("cfg_link_path").Parse(cfg.LinkPath)
	if err != nil {
		return err
	}
	cfg.LinkPathTpl = tpl

	for pkgName, pkg := range cfg.Packages {
		if pkg.Result == nil {
			pkg.Result = &domain.PackageResult{
				Name:  pkgName,
				Files: map[string]domain.FileResult{},
			}
		}

		if pkg.LinkPath == "" {
			pkg.LinkPath = cfg.LinkPath
			pkg.LinkPathTpl = cfg.LinkPathTpl
		} else {
			tpl, err := template.New("pkg_link_path").Parse(pkg.LinkPath)
			if err != nil {
				return err
			}
			pkg.LinkPathTpl = tpl
		}

		if pkg.BinPath == "" {
			pkg.BinPath = cfg.BinPath
			pkg.BinPathTpl = cfg.BinPathTpl
		} else {
			tpl, err := template.New("pkg_bin_path").Parse(pkg.BinPath)
			if err != nil {
				return err
			}
			pkg.BinPathTpl = tpl
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
		for i, file := range pkg.Files {
			if file.Result == nil {
				file.Result = &domain.FileResult{}
			}
			if file.Mode == 0 {
				file.Mode = 0755
			}

			if file.LinkPath == "" {
				file.LinkPath = pkg.LinkPath
				file.LinkPathTpl = pkg.LinkPathTpl
			} else {
				tpl, err := template.New("file_link_path").Parse(file.LinkPath)
				if err != nil {
					return err
				}
				file.LinkPathTpl = tpl
			}

			if file.BinPath == "" {
				file.BinPath = pkg.BinPath
				file.BinPathTpl = pkg.BinPathTpl
			} else {
				tpl, err := template.New("file_bin_path").Parse(file.BinPath)
				if err != nil {
					return err
				}
				file.BinPathTpl = tpl
			}

			dst, err := util.RenderTpl(
				file.BinPathTpl, &domain.TemplateParams{
					Name: file.Name, Version: pkg.Version,
				})
			if err != nil {
				return err
			}
			file.Bin = dst

			lnPath, err := util.RenderTpl(
				file.LinkPathTpl, &domain.TemplateParams{
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
