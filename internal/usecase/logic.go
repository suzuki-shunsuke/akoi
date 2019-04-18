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

func NewLogic(param domain.LogicParam) domain.Logic {
	lgc := &logic{
		fsys:  param.Fsys,
		logic: param.Logic,
	}
	if lgc.logic == nil {
		lgc.logic = lgc
	}
	return lgc
}
