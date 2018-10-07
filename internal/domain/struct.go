package domain

import (
	"net/url"
	"os"
	"text/template"
)

type (
	// Config represents application's configuration.
	Config struct {
		BinPathTpl        *template.Template `yaml:"-"`
		LinkPathTpl       *template.Template `yaml:"-"`
		BinPath           string             `yaml:"bin_path"`
		LinkPath          string             `yaml:"link_path"`
		NumOfDLPartitions int                `yaml:"num_of_dl_partitions"`
		Packages          map[string]Package `yaml:"packages"`
	}

	// File represents a file configuration.
	File struct {
		Name        string             `yaml:"name"`
		Archive     string             `yaml:"archive"`
		Bin         string             `yaml:"-"`
		Link        string             `yaml:"-"`
		BinPath     string             `yaml:"bin_path"`
		LinkPath    string             `yaml:"link_path"`
		BinPathTpl  *template.Template `yaml:"-"`
		LinkPathTpl *template.Template `yaml:"-"`
		Mode        os.FileMode        `yaml:"mode"`
		Result      *FileResult        `yaml:"-"`
	}

	// FileResult represents a result of file installation.
	FileResult struct {
		Error       string `json:"error"`
		FileRemoved bool   `json:"file_removed"`
		Changed     bool   `json:"changed"`
		Migrated    bool   `json:"migrated"`
		ModeChanged bool   `json:"mode_changed"`
		Installed   bool   `json:"installed"`
		DirCreated  bool   `json:"dir_created"`
		Name        string `json:"name"`
		Link        string `json:"link"`
		Entity      string `json:"entity"`
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
		Chmod           Chmod           `validate:"required"`
		Copy            Copy            `validate:"required"`
		Download        Download        `validate:"required"`
		ExpandEnv       ExpandEnv       `validate:"required"`
		Fprintf         Fprintf         `validate:"required"`
		Fprintln        Fprintln        `validate:"required"`
		GetArchiver     GetArchiver     `validate:"required"`
		GetFileStat     GetFileStat     `validate:"required"`
		GetFileLstat    GetFileStat     `validate:"required"`
		MkdirAll        MkdirAll        `validate:"required"`
		MkLink          MkLink          `validate:"required"`
		NewGzipReader   NewGzipReader   `validate:"required"`
		NewLoggerOutput NewLoggerOutput `validate:"required"`
		Open            Open            `validate:"required"`
		OpenFile        OpenFile        `validate:"required"`
		Printf          Printf          `validate:"required"`
		Println         Println         `validate:"required"`
		ReadConfigFile  ReadConfigFile  `validate:"required"`
		ReadLink        ReadLink        `validate:"required"`
		RemoveAll       RemoveFile      `validate:"required"`
		RemoveFile      RemoveFile      `validate:"required"`
		RemoveLink      RemoveFile      `validate:"required"`
		TempDir         TempDir         `validate:"required"`
	}

	// InstallParams is parameters of usecase.Install .
	InstallParams struct {
		ConfigFilePath string
		Format         string
		DryRun         bool
	}

	// Package represents a package configuration.
	Package struct {
		ArchiveType       string             `yaml:"archive_type"`
		Name              string             `yaml:"-" validate:"required"`
		RawURL            string             `yaml:"url" validate:"required"`
		Version           string             `yaml:"version" validate:"required"`
		BinPath           string             `yaml:"bin_path"`
		LinkPath          string             `yaml:"link_path"`
		NumOfDLPartitions int                `yaml:"num_of_dl_partitions"`
		BinPathTpl        *template.Template `yaml:"-"`
		LinkPathTpl       *template.Template `yaml:"-"`
		Archiver          Archiver           `yaml:"-" validate:"required"`
		Files             []File             `yaml:"files"`
		URL               *url.URL           `yaml:"-"`
		Result            *PackageResult     `yaml:"-"`
	}

	// PackageResult represents a result of package installation.
	PackageResult struct {
		Error   string                `json:"error"`
		Name    string                `json:"-"`
		Version string                `json:"version"`
		URL     string                `json:"url"`
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
