package usecase

import (
	"fmt"
	"net/url"
	"path/filepath"
	"text/template"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
	"github.com/suzuki-shunsuke/akoi/internal/util"
)

func (lgc *Logic) SetupConfig(cfg domain.Config) (domain.Config, error) {
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

	numCPUs := lgc.Runtime.NumCPU()
	if cfg.NumOfDLPartitions <= 0 {
		cfg.NumOfDLPartitions = numCPUs
	}

	for pkgName, pkg := range cfg.Packages {
		pkg, err := lgc.Logic.SetupPkgConfig(cfg, pkgName, pkg, numCPUs)
		if err != nil {
			return cfg, err
		}
		cfg.Packages[pkgName] = pkg
	}
	return cfg, nil
}

func (lgc *Logic) SetupPkgConfig(
	cfg domain.Config, pkgName string, pkg domain.Package, numCPUs int,
) (domain.Package, error) {
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
			return pkg, err
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
			return pkg, err
		}
		pkg.BinPathTpl = tpl
	}

	pkg.Name = pkgName
	tpl, err := template.New("pkg_url").Parse(pkg.RawURL)
	if err != nil {
		return pkg, err
	}
	u, err := util.RenderTpl(tpl, lgc.getTemplateParams(&pkg, nil))
	if err != nil {
		return pkg, err
	}
	u2, err := url.Parse(u)
	if err != nil {
		return pkg, err
	}
	pkg.URL = u2
	pkg.Archiver = lgc.GetArchiver.Get(u2.Path, pkg.ArchiveType)

	if pkg.NumOfDLPartitions < 0 {
		pkg.NumOfDLPartitions = numCPUs
	} else {
		if pkg.NumOfDLPartitions == 0 {
			pkg.NumOfDLPartitions = cfg.NumOfDLPartitions
		}
	}

	for i, file := range pkg.Files {
		file, err := lgc.Logic.SetupFileConfig(pkg, file)
		if err != nil {
			return pkg, err
		}
		pkg.Files[i] = file
	}
	return pkg, nil
}

func (lgc *Logic) SetupFileConfig(
	pkg domain.Package, file domain.File,
) (domain.File, error) {
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
			return file, err
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
			return file, err
		}
		file.BinPathTpl = tpl
	}

	tplParams := lgc.getTemplateParams(&pkg, &file)

	dst, err := util.RenderTpl(file.BinPathTpl, tplParams)
	if err != nil {
		return file, err
	}
	if !filepath.IsAbs(dst) {
		return file, fmt.Errorf(
			"installed path must be absolute: %s %s %s", pkg.Name, file.Name, dst)
	}
	file.Bin = dst

	arcPath := lgc.Fsys.ExpandEnv(file.Archive)
	arcPathTpl, err := template.New("archive_path").Parse(arcPath)
	if err != nil {
		return file, err
	}
	file.Archive, err = util.RenderTpl(arcPathTpl, tplParams)
	if err != nil {
		return file, err
	}

	lnPath, err := util.RenderTpl(
		file.LinkPathTpl, tplParams)
	if err != nil {
		return file, err
	}
	if !filepath.IsAbs(lnPath) {
		return file, fmt.Errorf(
			"link path must be absolute: %s %s %s", pkg.Name, file.Name, lnPath)
	}
	file.Link = lnPath
	return file, nil
}

func (lgc *Logic) getTemplateParams(pkg *domain.Package, file *domain.File) *domain.TemplateParams {
	params := &domain.TemplateParams{
		Version: pkg.Version,
		OS:      lgc.Runtime.OS(),
		Arch:    lgc.Runtime.Arch(),
	}
	if file != nil {
		params.Name = file.Name
	}
	return params
}
