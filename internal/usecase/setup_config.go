package usecase

import (
	"net/url"
	"path/filepath"
	"text/template"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
	"github.com/suzuki-shunsuke/akoi/internal/util"
)

func (lgc *Logic) SetupConfig(cfg domain.Config, cfgPath string) (domain.Config, error) {
	cfgDir, err := filepath.Abs(lgc.Fsys.ExpandEnv(filepath.Dir(cfgPath)))
	if err != nil {
		return cfg, err
	}
	cfg.BinPathTpl, err = lgc.parseCfgBinAndLinkPath(
		cfg.BinPath, cfgDir, "cfg_bin_path")
	if err != nil {
		return cfg, err
	}
	cfg.LinkPathTpl, err = lgc.parseCfgBinAndLinkPath(
		cfg.LinkPath, cfgDir, "cfg_link_path")
	if err != nil {
		return cfg, err
	}

	numCPUs := lgc.Runtime.NumCPU()
	if cfg.NumOfDLPartitions <= 0 {
		cfg.NumOfDLPartitions = numCPUs
	}

	for pkgName, pkg := range cfg.Packages {
		pkg, err := lgc.Logic.SetupPkgConfig(cfg, cfgDir, pkgName, pkg, numCPUs)
		if err != nil {
			return cfg, err
		}
		cfg.Packages[pkgName] = pkg
	}
	return cfg, nil
}

func (lgc *Logic) SetupPkgConfig(
	cfg domain.Config, cfgDir, pkgName string, pkg domain.Package, numCPUs int,
) (domain.Package, error) {
	if pkg.Result == nil {
		pkg.Result = &domain.PackageResult{
			Name:  pkgName,
			Files: map[string]domain.FileResult{},
		}
	}

	var err error
	if pkg.LinkPath == "" {
		pkg.LinkPathTpl = cfg.LinkPathTpl
	} else {
		pkg.LinkPathTpl, err = lgc.parseCfgBinAndLinkPath(
			pkg.LinkPath, cfgDir, "pkg_link_path")
		if err != nil {
			return pkg, err
		}
	}
	if pkg.BinPath == "" {
		pkg.BinPathTpl = cfg.BinPathTpl
	} else {
		pkg.BinPathTpl, err = lgc.parseCfgBinAndLinkPath(
			pkg.BinPath, cfgDir, "pkg_bin_path")
		if err != nil {
			return pkg, err
		}
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

	if pkg.HTTPRequestTimeout == 0 {
		pkg.HTTPRequestTimeout = cfg.HTTPRequestTimeout
	}

	for i, file := range pkg.Files {
		file, err := lgc.Logic.SetupFileConfig(pkg, cfgDir, file)
		if err != nil {
			return pkg, err
		}
		pkg.Files[i] = file
	}
	return pkg, nil
}

func (lgc *Logic) SetupFileConfig(
	pkg domain.Package, cfgDir string, file domain.File,
) (domain.File, error) {
	if file.Result == nil {
		file.Result = &domain.FileResult{}
	}
	if file.Mode == 0 {
		file.Mode = 0755
	}

	var err error
	if file.LinkPath == "" {
		file.LinkPathTpl = pkg.LinkPathTpl
	} else {
		file.LinkPathTpl, err = lgc.parseCfgBinAndLinkPath(
			file.LinkPath, cfgDir, "file_link_path")
		if err != nil {
			return file, err
		}
	}
	if file.BinPath == "" {
		file.BinPathTpl = pkg.BinPathTpl
	} else {
		file.BinPathTpl, err = lgc.parseCfgBinAndLinkPath(
			file.BinPath, cfgDir, "file_bin_path")
		if err != nil {
			return file, err
		}
	}

	tplParams := lgc.getTemplateParams(&pkg, &file)

	file.Bin, err = util.RenderTpl(file.BinPathTpl, tplParams)
	if err != nil {
		return file, err
	}

	arcPath := lgc.Fsys.ExpandEnv(file.Archive)
	arcPathTpl, err := template.New("archive_path").Parse(arcPath)
	if err != nil {
		return file, err
	}
	file.Archive, err = util.RenderTpl(arcPathTpl, tplParams)
	if err != nil {
		return file, err
	}

	file.Link, err = util.RenderTpl(file.LinkPathTpl, tplParams)
	return file, err
}

func (lgc *Logic) getTemplateParams(
	pkg *domain.Package, file *domain.File,
) *domain.TemplateParams {
	params := &domain.TemplateParams{
		Version: pkg.Version,
		OS:      lgc.Runtime.OS(),
		Arch:    lgc.Runtime.Arch(),
	}
	if file != nil {
		params.Name = file.Name
	} else {
		params.Name = pkg.Name
	}
	return params
}

func (lgc *Logic) parseCfgBinAndLinkPath(
	p, cfgDir, tplName string,
) (*template.Template, error) {
	if p == "" {
		return nil, nil
	}
	p = lgc.Fsys.ExpandEnv(p)
	if !filepath.IsAbs(p) {
		p = filepath.Join(cfgDir, p)
	}
	return template.New(tplName).Parse(p)
}
