package registry

import (
	"io"
	"os"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
	"github.com/suzuki-shunsuke/akoi/internal/infra"
	"github.com/suzuki-shunsuke/akoi/internal/testutil"
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
func NewInstallMethods(params domain.InstallParams) domain.InstallMethods {
	flag := params.Format != "ansible"
	if params.DryRun {
		return domain.InstallMethods{
			Chmod: testutil.NewFakeChmod(nil),
			Copy:  testutil.NewFakeCopy(10, nil),
			Download: testutil.NewFakeDownload(
				testutil.NewFakeIOReadCloser("hello"), nil),
			ExpandEnv:    os.ExpandEnv,
			Fprintf:      infra.NewFprintf(flag),
			Fprintln:     infra.NewFprintln(flag),
			GetArchiver:  testutil.NewFakeGetArchiver(nil),
			GetFileStat:  os.Stat,
			GetFileLstat: os.Lstat,
			MkdirAll:     testutil.NewFakeMkdirAll(nil),
			MkLink:       testutil.NewFakeMkLink(nil),
			NewGzipReader: testutil.NewFakeNewGzipReader(
				testutil.NewFakeIOReadCloser("hello"), nil),
			NewLoggerOutput: infra.NewLoggerOutput,
			Open:            testutil.NewFakeOpen(&os.File{}, nil),
			OpenFile:        testutil.NewFakeOpenFile(&os.File{}, nil),
			Printf:          infra.NewPrintf(flag),
			Println:         infra.NewPrintln(flag),
			ReadConfigFile:  infra.ReadConfigFile,
			ReadLink:        os.Readlink,
			RemoveAll:       testutil.NewFakeRemoveFile(nil),
			RemoveLink:      testutil.NewFakeRemoveFile(nil),
			RemoveFile:      testutil.NewFakeRemoveFile(nil),
			TempDir:         testutil.NewFakeTempDir("/tmp/tempdir", nil),
		}
	}
	return domain.InstallMethods{
		Chmod:           os.Chmod,
		Copy:            io.Copy,
		Download:        infra.Download,
		ExpandEnv:       os.ExpandEnv,
		Fprintf:         infra.NewFprintf(flag),
		Fprintln:        infra.NewFprintln(flag),
		GetArchiver:     infra.GetArchiver,
		GetFileStat:     os.Stat,
		GetFileLstat:    os.Lstat,
		MkdirAll:        infra.MkdirAll,
		MkLink:          os.Symlink,
		NewGzipReader:   infra.NewGzipReader,
		NewLoggerOutput: infra.NewLoggerOutput,
		Open:            os.Open,
		OpenFile:        os.OpenFile,
		Printf:          infra.NewPrintf(flag),
		Println:         infra.NewPrintln(flag),
		ReadConfigFile:  infra.ReadConfigFile,
		ReadLink:        os.Readlink,
		RemoveAll:       os.RemoveAll,
		RemoveLink:      os.Remove,
		RemoveFile:      os.Remove,
		TempDir:         infra.TempDir,
	}
}
