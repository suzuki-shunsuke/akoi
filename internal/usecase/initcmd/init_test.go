package initcmd

import (
	"fmt"
	"testing"

	"github.com/suzuki-shunsuke/gomic/gomic"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
	"github.com/suzuki-shunsuke/akoi/internal/test"
)

func TestInitConfigFile(t *testing.T) {
	fsys := test.NewFileSystem(t, gomic.DoNothing).
		SetReturnExistFile(true)
	params := &domain.InitParams{Dest: "dest"}
	if err := InitConfigFile(params, fsys); err != nil {
		t.Fatal(err)
	}
	fsys.SetReturnExistFile(false)
	if err := InitConfigFile(params, fsys); err != nil {
		t.Fatal(err)
	}
	fsys.SetReturnMkdirAll(fmt.Errorf("failed to create a directory"))
	if err := InitConfigFile(params, fsys); err == nil {
		t.Fatal("it should be failed to create a directory")
	}
}
