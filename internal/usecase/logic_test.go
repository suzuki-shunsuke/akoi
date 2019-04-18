package usecase

import (
	"testing"

	"github.com/suzuki-shunsuke/gomic/gomic"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
	"github.com/suzuki-shunsuke/akoi/internal/test"
)

func newLogicParam(t *testing.T) domain.LogicParam {
	return domain.LogicParam{
		Fsys: test.NewFileSystem(t, gomic.DoNothing),
	}
}

func TestNewLogic(t *testing.T) {
	lgc := NewLogic(newLogicParam(t))
	if lgc == nil {
		t.Fatal("logic is nil")
	}
}
