package usecase

import (
	"testing"

	"github.com/suzuki-shunsuke/gomic/gomic"

	"github.com/suzuki-shunsuke/akoi/internal/test"
)

func newLogicMock(t *testing.T) *Logic {
	return &Logic{
		Logic:   test.NewLogic(t, gomic.DoNothing),
		Fsys:    test.NewFileSystem(t, gomic.DoNothing),
		Printer: test.NewPrinter(t, gomic.DoNothing),
	}
}
