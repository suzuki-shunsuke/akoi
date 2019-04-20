package usecase

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suzuki-shunsuke/gomic/gomic"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
	"github.com/suzuki-shunsuke/akoi/internal/test"
)

func TestLogicGetInstalledFiles(t *testing.T) {
	data := []struct {
		title string
		files []domain.File
		exp   []domain.File
		fsys  domain.FileSystem
	}{
		{
			title: "files is empty",
			exp:   []domain.File{}},
		{
			title: "no file is installed",
			files: []domain.File{{}},
			exp:   []domain.File{},
			fsys: test.NewFileSystem(t, gomic.DoNothing).
				SetReturnGetFileStat(test.NewFileInfo(t, gomic.DoNothing), nil),
		}, {
			title: "failed to change a file permission",
			files: []domain.File{{Mode: 0644, Result: &domain.FileResult{}}},
			exp:   []domain.File{},
			fsys: test.NewFileSystem(t, gomic.DoNothing).
				SetReturnChmod(fmt.Errorf("permission denied")).
				SetReturnGetFileStat(test.NewFileInfo(t, gomic.DoNothing), nil),
		}, {
			title: "change a file permission",
			files: []domain.File{{Mode: 0644, Result: &domain.FileResult{}}},
			exp:   []domain.File{},
			fsys: test.NewFileSystem(t, gomic.DoNothing).
				SetReturnGetFileStat(test.NewFileInfo(t, gomic.DoNothing), nil),
		}, {
			title: "failed to create a parent directory",
			files: []domain.File{{Result: &domain.FileResult{}}},
			exp:   []domain.File{},
			fsys: test.NewFileSystem(t, gomic.DoNothing).
				SetReturnMkdirAll(fmt.Errorf("permission denied")).
				SetReturnGetFileStat(nil, fmt.Errorf("file is not found")),
		}, {
			title: "create a parent directory",
			files: []domain.File{{Result: &domain.FileResult{}}},
			exp: []domain.File{{
				Result: &domain.FileResult{
					DirCreated: true,
				},
			}},
			fsys: test.NewFileSystem(t, gomic.DoNothing).
				SetReturnGetFileStat(nil, fmt.Errorf("file is not found")),
		},
	}
	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			logic := newLogicMock(t)
			if d.fsys != nil {
				logic.Fsys = d.fsys
			}
			require.Equal(t, d.exp, logic.GetInstalledFiles(d.files))
		})
	}
}
