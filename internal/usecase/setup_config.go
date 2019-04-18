package usecase

import (
	"fmt"
	"net/url"
	"path/filepath"
	"runtime"
	"text/template"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
	"github.com/suzuki-shunsuke/akoi/internal/util"
)

func (lgc *Logic) SetupConfig(
	cfg domain.Config, getArchiver domain.GetArchiver,
) (domain.Config, error) {
	cfg.BinPath = lgc.Fsys.ExpandEnv(cfg.BinPath)
	cfg.LinkPath = lgc.Fsys.ExpandEnv(cfg.LinkPath)
	tpl, err := template.New("cfg_bin_path").Parse(cfg.BinPath)
	if err != nil {
		return cfg, err
	}
	cfg.BinPathTpl = tpl

	tpl, err = template.New("cfg_link_path").Parse(cfg.LinkPath)
	if err != nil {
		return cfg, err
	}
	cfg.LinkPathTpl = tpl

	numCPUs := runtime.NumCPU()
	if cfg.NumOfDLPartitions <= 0 {
		cfg.NumOfDLPartitions = numCPUs
	}

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
			pkg.LinkPath = lgc.Fsys.ExpandEnv(pkg.LinkPath)
			tpl, err := template.New("pkg_link_path").Parse(pkg.LinkPath)
			if err != nil {
				return cfg, err
			}
			pkg.LinkPathTpl = tpl
		}

		if pkg.BinPath == "" {
			pkg.BinPath = cfg.BinPath
			pkg.BinPathTpl = cfg.BinPathTpl
		} else {
			pkg.BinPath = lgc.Fsys.ExpandEnv(pkg.BinPath)
			tpl, err := template.New("pkg_bin_path").Parse(pkg.BinPath)
			if err != nil {
				return cfg, err
			}
			pkg.BinPathTpl = tpl
		}

		pkg.Name = pkgName
		tpl, err := template.New("pkg_url").Parse(pkg.RawURL)
		if err != nil {
			return cfg, err
		}
		u, err := util.RenderTpl(tpl, pkg)
		if err != nil {
			return cfg, err
		}
		u2, err := url.Parse(u)
		if err != nil {
			return cfg, err
		}
		pkg.URL = u2
		pkg.Archiver = getArchiver.Get(u2.Path, pkg.ArchiveType)

		if pkg.NumOfDLPartitions < 0 {
			pkg.NumOfDLPartitions = numCPUs
		} else {
			if pkg.NumOfDLPartitions == 0 {
				pkg.NumOfDLPartitions = cfg.NumOfDLPartitions
			}
		}

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
				file.LinkPath = lgc.Fsys.ExpandEnv(file.LinkPath)
				tpl, err := template.New("file_link_path").Parse(file.LinkPath)
				if err != nil {
					return cfg, err
				}
				file.LinkPathTpl = tpl
			}

			if file.BinPath == "" {
				file.BinPath = pkg.BinPath
				file.BinPathTpl = pkg.BinPathTpl
			} else {
				file.BinPath = lgc.Fsys.ExpandEnv(file.BinPath)
				tpl, err := template.New("file_bin_path").Parse(file.BinPath)
				if err != nil {
					return cfg, err
				}
				file.BinPathTpl = tpl
			}

			dst, err := util.RenderTpl(
				file.BinPathTpl, &domain.TemplateParams{
					Name: file.Name, Version: pkg.Version,
				})
			if err != nil {
				return cfg, err
			}
			if !filepath.IsAbs(dst) {
				return cfg, fmt.Errorf(
					"installed path must be absolute: %s %s %s", pkg.Name, file.Name, dst)
			}
			file.Bin = dst

			arcPath := lgc.Fsys.ExpandEnv(file.Archive)
			arcPathTpl, err := template.New("archive_path").Parse(arcPath)
			if err != nil {
				return cfg, err
			}
			file.Archive, err = util.RenderTpl(
				arcPathTpl, &domain.TemplateParams{
					Name: file.Name, Version: pkg.Version,
				})
			if err != nil {
				return cfg, err
			}

			lnPath, err := util.RenderTpl(
				file.LinkPathTpl, &domain.TemplateParams{
					Name: file.Name, Version: pkg.Version,
				})
			if err != nil {
				return cfg, err
			}
			if !filepath.IsAbs(lnPath) {
				return cfg, fmt.Errorf(
					"link path must be absolute: %s %s %s", pkg.Name, file.Name, lnPath)
			}
			file.Link = lnPath
			pkg.Files[i] = file
		}
		cfg.Packages[pkgName] = pkg
	}
	return cfg, nil
}
