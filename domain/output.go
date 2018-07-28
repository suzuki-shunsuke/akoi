package domain

import (
	"encoding/json"
)

type (
	// FileResult represents a result of file installation.
	FileResult struct {
		Error   string `json:"error"`
		Changed bool   `json:"changed"`
		// installed, migrated, changed, none
		State  string `json:"state"`
		Name   string `json:"name"`
		Link   string `json:"link"`
		Entity string `json:"entity"`
	}

	// PackageResult represents a result of package installation.
	PackageResult struct {
		Error   string       `json:"error"`
		Name    string       `json:"-"`
		State   string       `json:"state"`
		Version string       `json:"version"`
		URL     string       `json:"url"`
		Changed bool         `json:"changed"`
		Files   []FileResult `json:"files"`
	}

	// Result represents a result of packages's installation.
	Result struct {
		Msg      string                   `json:"msg"`
		Changed  bool                     `json:"changed"`
		Failed   bool                     `json:"failed"`
		Packages map[string]PackageResult `json:"packages"`
	}
)

// String converts result into a string.
func (result *Result) String(params *InstallParams) string {
	switch params.Format {
	case "ansible":
		b, err := json.Marshal(result)
		if err != nil {
			return `{"changed": true}`
		}
		return string(b)
	default:
		return ""
	}
}
