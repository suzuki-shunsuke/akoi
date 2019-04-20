package usecase

import (
	"fmt"
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
		},
	}
	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			logic := newLogicMock(t)
			logic.Fsys = d.fsys
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
