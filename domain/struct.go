package domain

import (
	"net/url"
	"os"
	"text/template"
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
		Name     string   `yaml:"-" validate:"required"`
		URL      *url.URL `yaml:"-"`
		Archiver Archiver `yaml:"-" validate:"required"`
		RawURL   string   `yaml:"url" validate:"required"`
		Version  string   `yaml:"version" validate:"required"`
		Files    []File   `yaml:"files"`
	}

	// TemplateParams is template parameters.
	TemplateParams struct {
		Name    string
		Version string
	}
)
