package usecase

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suzuki-shunsuke/gomic/gomic"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
	"github.com/suzuki-shunsuke/akoi/internal/test"
)

func TestLogicInstallFile(t *testing.T) {
	data := []struct {
		title         string
		fsys          domain.FileSystem
		getGzipReader domain.GetGzipReader
		pkg           domain.Package
		isErr         bool
	}{
		{
			title: "failed to open a dest file",
			fsys: test.NewFileSystem(t, gomic.DoNothing).
				SetReturnOpenFile(nil, fmt.Errorf("file is not found")),
			isErr: true,
		}, {
			title: "failed to open a source file",
			fsys: test.NewFileSystem(t, gomic.DoNothing).
				SetReturnOpen(nil, fmt.Errorf("file is not found")).
				SetReturnOpenFile(test.NewWriteCloser(t, gomic.DoNothing), nil),
			isErr: true,
		}, {
			title: "failed to copy a source file",
			fsys: test.NewFileSystem(t, gomic.DoNothing).
				SetReturnOpen(ioutil.NopCloser(bytes.NewBufferString("")), nil).
				SetReturnCopy(0, fmt.Errorf("permission denied")).
				SetReturnOpenFile(test.NewWriteCloser(t, gomic.DoNothing), nil),
			isErr: true,
		}, {
			title: "copy a source file",
			fsys: test.NewFileSystem(t, gomic.DoNothing).
				SetReturnOpen(ioutil.NopCloser(bytes.NewBufferString("")), nil).
				SetReturnOpenFile(test.NewWriteCloser(t, gomic.DoNothing), nil),
			isErr: false,
		}, {
			title: "failed to get a gzip reader",
			pkg:   domain.Package{ArchiveType: "Gzip"},
			fsys: test.NewFileSystem(t, gomic.DoNothing).
				SetReturnOpenFile(test.NewWriteCloser(t, gomic.DoNothing), nil),
			getGzipReader: test.NewGetGzipReader(t, gomic.DoNothing).
				SetReturnGet(nil, fmt.Errorf("permission denied")),
			isErr: true,
		}, {
			title: "failed to copy a gzip file",
			pkg:   domain.Package{ArchiveType: "Gzip"},
			getGzipReader: test.NewGetGzipReader(t, gomic.DoNothing).
				SetReturnGet(ioutil.NopCloser(bytes.NewBufferString("")), nil),
			fsys: test.NewFileSystem(t, gomic.DoNothing).
				SetReturnOpenFile(test.NewWriteCloser(t, gomic.DoNothing), nil).
				SetReturnCopy(0, fmt.Errorf("permission denied")),
			isErr: true,
		},
	}
	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			logic := newLogicMock(t)
			if d.fsys != nil {
				logic.Fsys = d.fsys
			}
			if d.getGzipReader != nil {
				logic.GetGzipReader = d.getGzipReader
			}
			err := logic.InstallFile(
				&domain.File{Result: &domain.FileResult{}}, d.pkg, domain.InstallParams{}, "", nil)
			if d.isErr {
				require.NotNil(t, err)
				return
			}
			require.Nil(t, err)
		})
	}
}
