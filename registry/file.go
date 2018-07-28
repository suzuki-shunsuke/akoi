package registry

import (
	"net/http"
	"os"

	"github.com/suzuki-shunsuke/akoi/domain"
	"github.com/suzuki-shunsuke/akoi/infra"
	"github.com/suzuki-shunsuke/akoi/testutil"
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
func NewInstallMethods(dryRun bool) *domain.InstallMethods {
	if dryRun {
		return &domain.InstallMethods{
			Chmod:    testutil.NewFakeChmod(nil),
			CopyFile: testutil.NewFakeCopyFile(nil),
			Download: testutil.NewFakeDownload(
				&http.Response{
					StatusCode: 200,
					Body:       testutil.NewFakeIOReadCloser("hello"),
				}, nil),
			Exist:          infra.ExistFile,
			GetArchiver:    testutil.NewFakeGetArchiver(nil),
			GetFileStat:    os.Stat,
			GetFileLstat:   os.Lstat,
			MkdirAll:       testutil.NewFakeMkdirAll(nil),
			MkLink:         testutil.NewFakeMkLink(nil),
			ReadConfigFile: infra.ReadConfigFile,
			ReadLink:       os.Readlink,
			RemoveAll:      testutil.NewFakeRemoveAll(nil),
			RemoveLink:     testutil.NewFakeRemoveLink(nil),
			RemoveFile:     testutil.NewFakeRemoveFile(nil),
			TempDir:        testutil.NewFakeTempDir("/tmp/tempdir", nil),
		}
	}
	return &domain.InstallMethods{
		Chmod:          os.Chmod,
		CopyFile:       infra.CopyFile,
		Download:       http.Get,
		Exist:          infra.ExistFile,
		GetArchiver:    infra.GetArchiver,
		GetFileStat:    os.Stat,
		GetFileLstat:   os.Lstat,
		MkdirAll:       infra.MkdirAll,
		MkLink:         os.Symlink,
		ReadConfigFile: infra.ReadConfigFile,
		ReadLink:       os.Readlink,
		RemoveAll:      os.RemoveAll,
		RemoveLink:     os.Remove,
		RemoveFile:     os.Remove,
		TempDir:        infra.TempDir,
	}
}
