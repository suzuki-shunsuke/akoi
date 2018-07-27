package domain

import (
	"os"
	"text/template"

	"github.com/suzuki-shunsuke/akoi/util"
)

type (
	// Config represents application's configuration.
	Config struct {
		BinPathTpl  *template.Template
		LinkPathTpl *template.Template
		BinPath     string             `yaml:"bin_path"`
		LinkPath    string             `yaml:"link_path"`
		Packages    map[string]Package `yaml:"packages"`
	}

	// File represents a file configuration.
	File struct {
		Name    string      `yaml:"name"`
		Archive string      `yaml:"archive"`
		Bin     string      `yaml:"-"`
		Link    string      `yaml:"-"`
		Mode    os.FileMode `yaml:"mode"`
	}

	// InitMethods is functions which are used at usecase.Init .
	InitMethods struct {
		Write    WriteFile `validate:"required"`
		Exist    ExistFile `validate:"required"`
		MkdirAll MkdirAll  `validate:"required"`
	}

	// InitParams is parameters of usecase.Init .
	InitParams struct {
		Dest string
	}

	// InstallMethods is functions which are used at usecase.Install .
	InstallMethods struct {
		Chmod          Chmod          `validate:"required"`
		CopyFile       CopyFile       `validate:"required"`
		Download       Download       `validate:"required"`
		Exist          ExistFile      `validate:"required"`
		GetArchiver    GetArchiver    `validate:"required"`
		GetFileStat    GetFileStat    `validate:"required"`
		GetFileLstat   GetFileStat    `validate:"required"`
		MkdirAll       MkdirAll       `validate:"required"`
		MkLink         MkLink         `validate:"required"`
		ReadConfigFile ReadConfigFile `validate:"required"`
		ReadLink       ReadLink       `validate:"required"`
		RemoveAll      RemoveAll      `validate:"required"`
		RemoveLink     RemoveLink     `validate:"required"`
		TempDir        TempDir        `validate:"required"`
	}

	// InstallParams is parameters of usecase.Install .
	InstallParams struct {
		ConfigFilePath string
		Format         string
		DryRun         bool
	}

	// Package represents a package configuration.
	Package struct {
		Name    string `yaml:"-" validate:"required"`
		url     string
		URL     string `yaml:"url" validate:"required"`
		Version string `yaml:"version" validate:"required"`
		Files   []File `yaml:"files"`
	}

	// TemplateParams is template parameters.
	TemplateParams struct {
		Name    string
		Version string
	}
)

// GetURL returns a download URL.
func (pkg *Package) GetURL() string {
	return pkg.url
}

// Setup compiles and renders templates.
func (cfg *Config) Setup() error {
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
		tpl, err := template.New("pkg_url").Parse(pkg.URL)
		if err != nil {
			return err
		}
		u, err := util.RenderTpl(tpl, pkg)
		if err != nil {
			return err
		}
		pkg.url = u
		for i, file := range pkg.Files {
			dst, err := util.RenderTpl(
				cfg.BinPathTpl, &TemplateParams{
					Name: file.Name, Version: pkg.Version,
				})
			if err != nil {
				return err
			}
			file.Bin = dst

			lnPath, err := util.RenderTpl(
				cfg.LinkPathTpl, &TemplateParams{
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
