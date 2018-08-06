package registry

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

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
func NewInstallMethods(params *domain.InstallParams) *domain.InstallMethods {
	flag := params.Format != "ansible"
	if params.DryRun {
		return &domain.InstallMethods{
			Chmod: testutil.NewFakeChmod(nil),
			Copy:  testutil.NewFakeCopy(10, nil),
			Download: testutil.NewFakeDownload(
				&http.Response{
					StatusCode: 200,
					Body:       testutil.NewFakeIOReadCloser("hello"),
				}, nil),
			Fprintf:      infra.NewFprintf(flag),
			Fprintln:     infra.NewFprintln(flag),
			GetArchiver:  testutil.NewFakeGetArchiver(nil),
			GetFileStat:  os.Stat,
			GetFileLstat: os.Lstat,
			MkdirAll:     testutil.NewFakeMkdirAll(nil),
			MkLink:       testutil.NewFakeMkLink(nil),
			NewGzipReader: testutil.NewFakeNewGzipReader(
				testutil.NewFakeIOReadCloser("hello"), nil),
			Open:           testutil.NewFakeOpen(&os.File{}, nil),
			OpenFile:       testutil.NewFakeOpenFile(&os.File{}, nil),
			Printf:         infra.NewPrintf(flag),
			Println:        infra.NewPrintln(flag),
			ReadConfigFile: infra.ReadConfigFile,
			ReadLink:       os.Readlink,
			RemoveAll:      testutil.NewFakeRemoveFile(nil),
			RemoveLink:     testutil.NewFakeRemoveFile(nil),
			RemoveFile:     testutil.NewFakeRemoveFile(nil),
			TempDir:        testutil.NewFakeTempDir("/tmp/tempdir", nil),
		}
	}
	return &domain.InstallMethods{
		Chmod:          os.Chmod,
		Copy:           io.Copy,
		Download:       http.Get,
		Fprintf:        infra.NewFprintf(flag),
		Fprintln:       infra.NewFprintln(flag),
		GetArchiver:    infra.GetArchiver,
		GetFileStat:    os.Stat,
		GetFileLstat:   os.Lstat,
		MkdirAll:       infra.MkdirAll,
		MkLink:         os.Symlink,
		NewGzipReader:  infra.NewGzipReader,
		Open:           os.Open,
		OpenFile:       os.OpenFile,
		Printf:         infra.NewPrintf(flag),
		Println:        infra.NewPrintln(flag),
		ReadConfigFile: infra.ReadConfigFile,
		ReadLink:       os.Readlink,
		RemoveAll:      os.RemoveAll,
		RemoveLink:     os.Remove,
		RemoveFile:     os.Remove,
		TempDir:        infra.TempDir,
	}
}

// NewListMethods returns a domain.ListMethods .
func NewListMethods(params *domain.ListParams) *domain.ListMethods {
	flag := params.Format != "ansible"
	return &domain.ListMethods{
		Fprintf:        infra.NewFprintf(flag),
		Fprintln:       infra.NewFprintln(flag),
		GetArchiver:    infra.GetArchiver,
		Glob:           filepath.Glob,
		Printf:         infra.NewPrintf(flag),
		Println:        infra.NewPrintln(flag),
		ReadConfigFile: infra.ReadConfigFile,
	}
}
