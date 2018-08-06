package domain

import (
	"net/url"
	"os"
	"text/template"
)

type (
	// Config represents application's configuration.
	Config struct {
		BinDir       string             `yaml:"bin_dir"`
		BinSeparator string             `yaml:"bin_separator"`
		LinkDir      string             `yaml:"link_dir"`
		BinDirTpl    *template.Template `yaml:"-"`
		LinkDirTpl   *template.Template `yaml:"-"`
		Packages     map[string]Package `yaml:"packages"`
	}

	// File represents a file configuration.
	File struct {
		Archive      string             `yaml:"archive"`
		Bin          string             `yaml:"-"`
		BinDir       string             `yaml:"bin_dir"`
		BinSeparator string             `yaml:"bin_separator"`
		Link         string             `yaml:"-"`
		LinkDir      string             `yaml:"link_dir"`
		Name         string             `yaml:"name"`
		BinDirTpl    *template.Template `yaml:"-"`
		LinkDirTpl   *template.Template `yaml:"-"`
		Mode         os.FileMode        `yaml:"mode"`
		Result       *FileResult        `yaml:"-"`
	}

	// FileResult represents a result of file installation.
	FileResult struct {
		Entity      string `json:"entity"`
		Error       string `json:"error"`
		Link        string `json:"link"`
		Name        string `json:"name"`
		Changed     bool   `json:"changed"`
		DirCreated  bool   `json:"dir_created"`
		FileRemoved bool   `json:"file_removed"`
		Installed   bool   `json:"installed"`
		Migrated    bool   `json:"migrated"`
		ModeChanged bool   `json:"mode_changed"`
	}

	// InitMethods is functions which are used at usecase.Init .
	InitMethods struct {
		Exist    ExistFile `validate:"required"`
		MkdirAll MkdirAll  `validate:"required"`
		Write    WriteFile `validate:"required"`
	}

	// InitParams is parameters of usecase.Init .
	InitParams struct {
		Dest string
	}

	// InstallMethods is functions which are used at usecase.Install .
	InstallMethods struct {
		Chmod          Chmod          `validate:"required"`
		Copy           Copy           `validate:"required"`
		Download       Download       `validate:"required"`
		Fprintf        Fprintf        `validate:"required"`
		Fprintln       Fprintln       `validate:"required"`
		GetArchiver    GetArchiver    `validate:"required"`
		GetFileStat    GetFileStat    `validate:"required"`
		GetFileLstat   GetFileStat    `validate:"required"`
		MkdirAll       MkdirAll       `validate:"required"`
		MkLink         MkLink         `validate:"required"`
		NewGzipReader  NewGzipReader  `validate:"required"`
		Open           Open           `validate:"required"`
		OpenFile       OpenFile       `validate:"required"`
		Printf         Printf         `validate:"required"`
		Println        Println        `validate:"required"`
		ReadConfigFile ReadConfigFile `validate:"required"`
		ReadLink       ReadLink       `validate:"required"`
		RemoveAll      RemoveFile     `validate:"required"`
		RemoveFile     RemoveFile     `validate:"required"`
		RemoveLink     RemoveFile     `validate:"required"`
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
		ArchiveType  string             `yaml:"archive_type"`
		BinDir       string             `yaml:"bin_dir"`
		BinSeparator string             `yaml:"bin_separator"`
		LinkDir      string             `yaml:"link_dir"`
		Name         string             `yaml:"-" validate:"required"`
		RawURL       string             `yaml:"url" validate:"required"`
		Version      string             `yaml:"version" validate:"required"`
		BinDirTpl    *template.Template `yaml:"-"`
		LinkDirTpl   *template.Template `yaml:"-"`
		Archiver     Archiver           `yaml:"-" validate:"required"`
		Files        []File             `yaml:"files"`
		URL          *url.URL           `yaml:"-"`
		Result       *PackageResult     `yaml:"-"`
	}

	// PackageResult represents a result of package installation.
	PackageResult struct {
		Error   string                `json:"error"`
		Name    string                `json:"-"`
		URL     string                `json:"url"`
		Version string                `json:"version"`
		Changed bool                  `json:"changed"`
		Failed  bool                  `json:"failed"`
		Files   map[string]FileResult `json:"files"`
	}

	// Result represents a result of packages's installation.
	Result struct {
		Msg      string                   `json:"msg"`
		Changed  bool                     `json:"changed"`
		Failed   bool                     `json:"failed"`
		Packages map[string]PackageResult `json:"packages"`
	}

	// TemplateParams is template parameters.
	TemplateParams struct {
		Name    string
		Version string
	}
)

// Archived returns whether the package is archived.
func (pkg *Package) Archived() bool {
	return pkg.ArchiveType != "unarchived" && pkg.ArchiveType != "Gzip"
}
