package usecase

import (
	"github.com/suzuki-shunsuke/akoi/internal/domain"
)

type (
	logic struct {
		logic domain.Logic
		fsys  domain.FileSystem
	}
)

func NewLogic(fsys domain.FileSystem) domain.Logic {
	lgc := &logic{
		fsys: fsys,
	}
	lgc.logic = lgc
	return lgc
}
