package domain

import (
	"encoding/json"
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
		FileRemoved bool   `json:"file_removed"`
		Migrated    bool   `json:"migrated"`
		ModeChanged bool   `json:"mode_changed"`
		Installed   bool   `json:"installed"`
		DirCreated  bool   `json:"dir_created"`
		LinkCreated bool   `json:"link_created"`
		LinkRemoved bool   `json:"link_removed"`
		Error       string `json:"error"`
		Name        string `json:"name"`
		Link        string `json:"link"`
		Entity      string `json:"entity"`
	}

	fileResultJSON struct {
		FileRemoved bool   `json:"file_removed"`
		Migrated    bool   `json:"migrated"`
		ModeChanged bool   `json:"mode_changed"`
		Installed   bool   `json:"installed"`
		DirCreated  bool   `json:"dir_created"`
		LinkCreated bool   `json:"link_created"`
		LinkRemoved bool   `json:"link_removed"`
		Error       string `json:"error"`
		Name        string `json:"name"`
		Link        string `json:"link"`
		Entity      string `json:"entity"`

		Changed bool `json:"changed"`
		Failed  bool `json:"failed"`
	}

	// InitParams is parameters of usecase.Init .
	InitParams struct {
		Dest string
	}

	// InstallParams is parameters of usecase.Install .
	InstallParams struct {
		ConfigFilePath string
		Format         string
		DryRun         bool
	}

	// LogicParam is parameters of usecase.logic
	LogicParam struct {
		Logic Logic
		Fsys  FileSystem
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
		Files   map[string]FileResult `json:"files"`
	}

	packageResultJSON struct {
		Error   string                `json:"error"`
		Name    string                `json:"-"`
		Version string                `json:"version"`
		URL     string                `json:"url"`
		Files   map[string]FileResult `json:"files"`

		Changed bool `json:"changed"`
		Failed  bool `json:"failed"`
	}

	// Result represents a result of packages's installation.
	Result struct {
		Msg      string                   `json:"msg"`
		Packages map[string]PackageResult `json:"packages"`
	}

	resultJSON struct {
		Msg      string                   `json:"msg"`
		Packages map[string]PackageResult `json:"packages"`

		Changed bool `json:"changed"`
		Failed  bool `json:"failed"`
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

// Changed returns whether file has changed.
func (f *FileResult) Changed() bool {
	return f.FileRemoved || f.Migrated || f.ModeChanged || f.Installed || f.DirCreated || f.LinkCreated || f.LinkRemoved
}

// Failed returns whether file installation is failed.
func (f *FileResult) Failed() bool {
	return f.Error != ""
}

// MarshalJSON implements json.Marshaler interface.
func (f *FileResult) MarshalJSON() ([]byte, error) {
	return json.Marshal(&fileResultJSON{
		FileRemoved: f.FileRemoved,
		Migrated:    f.Migrated,
		ModeChanged: f.ModeChanged,
		Installed:   f.Installed,
		DirCreated:  f.DirCreated,
		LinkCreated: f.LinkCreated,
		LinkRemoved: f.LinkRemoved,
		Error:       f.Error,
		Name:        f.Name,
		Link:        f.Link,
		Entity:      f.Entity,

		Changed: f.Changed(),
		Failed:  f.Failed(),
	})
}

// Changed returns whether package has changed.
func (p *PackageResult) Changed() bool {
	for _, f := range p.Files {
		if f.Changed() {
			return true
		}
	}
	return false
}

// Failed returns whether package installation is failed.
func (p *PackageResult) Failed() bool {
	if p.Error != "" {
		return true
	}
	for _, f := range p.Files {
		if f.Failed() {
			return true
		}
	}
	return false
}

// MarshalJSON implements json.Marshaler interface.
func (p *PackageResult) MarshalJSON() ([]byte, error) {
	return json.Marshal(&packageResultJSON{
		Error:   p.Error,
		Name:    p.Name,
		Version: p.Version,
		URL:     p.URL,
		Files:   p.Files,

		Changed: p.Changed(),
		Failed:  p.Failed(),
	})
}

// Changed returns whether packages have changed.
func (r *Result) Changed() bool {
	for _, p := range r.Packages {
		if p.Changed() {
			return true
		}
	}
	return false
}

// Failed returns whether packages installation are failed.
func (r *Result) Failed() bool {
	if r.Msg != "" {
		return true
	}
	for _, p := range r.Packages {
		if p.Failed() {
			return true
		}
	}
	return false
}

// MarshalJSON implements json.Marshaler interface.
func (r *Result) MarshalJSON() ([]byte, error) {
	return json.Marshal(&resultJSON{
		Msg:      r.Msg,
		Packages: r.Packages,

		Changed: r.Changed(),
		Failed:  r.Failed(),
	})
}
