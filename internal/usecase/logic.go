package usecase

import (
	"github.com/suzuki-shunsuke/akoi/internal/domain"
)

type (
	logic struct {
		logic domain.Logic
	}
)

func NewLogic() domain.Logic {
	lgc := &logic{}
	lgc.logic = lgc
	return lgc
}
