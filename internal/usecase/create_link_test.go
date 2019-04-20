package usecase

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suzuki-shunsuke/gomic/gomic"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
	"github.com/suzuki-shunsuke/akoi/internal/test"
)

func TestLogicCreateLink(t *testing.T) {
	data := []struct {
		title string
		file  domain.File
		fsys  domain.FileSystem
		isErr bool
		exp   domain.File
	}{
		{
			title: "create link",
			file:  domain.File{Result: &domain.FileResult{}},
			fsys: test.NewFileSystem(t, gomic.DoNothing).
				SetReturnGetFileLstat(nil, fmt.Errorf("file is not found")),
			exp: domain.File{Result: &domain.FileResult{
				LinkCreated: true,
			}},
		},
		{
			title: "failed to create link",
			file:  domain.File{Result: &domain.FileResult{}},
			fsys: test.NewFileSystem(t, gomic.DoNothing).
				SetReturnGetFileLstat(nil, fmt.Errorf("file is not found")).
				SetReturnMkLink(fmt.Errorf("failed to create a link")),
			isErr: true,
		}, {
			title: "file is directory",
			file:  domain.File{Result: &domain.FileResult{}},
			fsys: test.NewFileSystem(t, gomic.DoNothing).
				SetReturnGetFileLstat(
					test.NewFileInfo(t, gomic.DoNothing).
						SetReturnMode(os.ModeDir), nil),
			isErr: true,
		}, {
			title: "file is a pipe",
			file:  domain.File{Result: &domain.FileResult{}},
			fsys: test.NewFileSystem(t, gomic.DoNothing).
				SetReturnGetFileLstat(
					test.NewFileInfo(t, gomic.DoNothing).
						SetReturnMode(os.ModeNamedPipe), nil),
			isErr: true,
		}, {
			title: "file is a regular file",
			file:  domain.File{Result: &domain.FileResult{}},
			fsys: test.NewFileSystem(t, gomic.DoNothing).
				SetReturnGetFileLstat(
					test.NewFileInfo(t, gomic.DoNothing), nil),
		}, {
			title: "file is a link",
			file:  domain.File{Result: &domain.FileResult{}},
			fsys: test.NewFileSystem(t, gomic.DoNothing).
				SetReturnGetFileLstat(
					test.NewFileInfo(t, gomic.DoNothing).
						SetReturnMode(os.ModeSymlink), nil),
		}, {
			title: "file is a link",
			file:  domain.File{Result: &domain.FileResult{}},
			fsys: test.NewFileSystem(t, gomic.DoNothing).
				SetReturnGetFileLstat(
					test.NewFileInfo(t, gomic.DoNothing).
						SetReturnMode(os.ModeIrregular), nil),
			isErr: true,
		},
	}
	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			logic := newLogicMock(t)
			if d.fsys != nil {
				logic.Fsys = d.fsys
			}
			file, err := logic.CreateLink(d.file)
			if d.isErr {
				require.NotNil(t, err)
				return
			}
			require.Nil(t, err)
			require.Equal(t, d.exp, file)
		})
	}
}

func TestLogicRemoveFileAndCreateLink(t *testing.T) {
	data := []struct {
		title string
		file  domain.File
		fsys  domain.FileSystem
		isErr bool
		exp   domain.File
	}{
		{
			title: "failed to remove a file",
			file:  domain.File{Result: &domain.FileResult{}},
			fsys: test.NewFileSystem(t, gomic.DoNothing).
				SetReturnRemoveFile(fmt.Errorf("permission denied")),
			isErr: true,
		}, {
			title: "failed to create a link",
			file:  domain.File{Result: &domain.FileResult{}},
			fsys: test.NewFileSystem(t, gomic.DoNothing).
				SetReturnMkLink(fmt.Errorf("failed to create a link")),
			isErr: true,
		}, {
			title: "normal",
			file:  domain.File{Result: &domain.FileResult{}},
			exp: domain.File{Result: &domain.FileResult{
				FileRemoved: true,
				Migrated:    true}},
		},
	}
	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			logic := newLogicMock(t)
			if d.fsys != nil {
				logic.Fsys = d.fsys
			}
			file, err := logic.RemoveFileAndCreateLink(d.file)
			if d.isErr {
				require.NotNil(t, err)
				return
			}
			require.Nil(t, err)
			require.Equal(t, d.exp, file)
		})
	}
}
