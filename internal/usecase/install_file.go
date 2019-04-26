package usecase

import (
	"io"
	"os"
	"path/filepath"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
)

func (lgc *Logic) InstallFile(
	file *domain.File, pkg domain.Package, params domain.InstallParams,
	tmpDir string, body io.Reader,
) error {
	lgc.Printer.Printf("install %s\n", file.Bin)
	writer, err := lgc.Fsys.OpenFile(file.Bin, os.O_RDWR|os.O_CREATE, file.Mode)
	if err != nil {
		return err
	}
	defer writer.Close()
	if pkg.Archived() {
		src, err := lgc.Fsys.Open(filepath.Join(tmpDir, file.Archive))
		if err != nil {
			return err
		}
		defer src.Close()
		if _, err := lgc.Fsys.Copy(writer, src); err != nil {
			return err
		}
		file.Result.Installed = true
		return nil
	}
	if pkg.ArchiveType == "Gzip" {
		reader, err := lgc.GetGzipReader.Get(body)
		if err != nil {
			return err
		}
		defer reader.Close()
		if _, err := lgc.Fsys.Copy(writer, reader); err != nil {
			return err
		}
		file.Result.Installed = true
		return nil
	}
	if _, err := lgc.Fsys.Copy(writer, body); err != nil {
		return err
	}
	file.Result.Installed = true
	return nil
}
