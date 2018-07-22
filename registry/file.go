package registry

import (
	"net/http"
	"os"

	"github.com/suzuki-shunsuke/akoi/domain"
	"github.com/suzuki-shunsuke/akoi/infra"
)

// NewInitMethods returns a domain.InitMethods .
func NewInitMethods() *domain.InitMethods {
	return &domain.InitMethods{
		Exist:    infra.ExistFile,
		Write:    infra.WriteFile,
		MkdirAll: infra.MkdirAll,
	}
}

// NewInstallMethods returns a domain.InstallMethods .
func NewInstallMethods() *domain.InstallMethods {
	return &domain.InstallMethods{
		Chmod:          os.Chmod,
		CopyFile:       infra.CopyFile,
		Download:       http.Get,
		Exist:          infra.ExistFile,
		GetFileStat:    os.Stat,
		GetFileLstat:   os.Lstat,
		MkdirAll:       infra.MkdirAll,
		MkLink:         os.Symlink,
		ReadConfigFile: infra.ReadConfigFile,
		ReadLink:       os.Readlink,
		RemoveAll:      os.RemoveAll,
		RemoveLink:     os.Remove,
		TempDir:        infra.TempDir,
	}
}
