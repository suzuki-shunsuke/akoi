package usecase

import (
	"github.com/suzuki-shunsuke/akoi/internal/domain"
)

type (
	Logic struct {
		Logic         domain.Logic
		Fsys          domain.FileSystem
		Printer       domain.Printer
		CfgReader     domain.ConfigReader
		Downloader    domain.Downloader
		GetArchiver   domain.GetArchiver
		GetGzipReader domain.GetGzipReader
	}
)
