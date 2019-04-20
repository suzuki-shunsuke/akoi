package initcmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suzuki-shunsuke/gomic/gomic"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
	"github.com/suzuki-shunsuke/akoi/internal/test"
)

func TestInitConfigFile(t *testing.T) {
	fsys := test.NewFileSystem(t, gomic.DoNothing).
		SetReturnExistFile(true)
	params := &domain.InitParams{Dest: "dest"}
	require.Nil(t, InitConfigFile(params, fsys))
	fsys.SetReturnExistFile(false)
	require.Nil(t, InitConfigFile(params, fsys))
	fsys.SetReturnMkdirAll(fmt.Errorf("failed to create a directory"))
	require.NotNil(t, InitConfigFile(params, fsys), "it should be failed to create a directory")
}
