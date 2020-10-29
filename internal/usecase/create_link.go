package usecase

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
)

func (lgc *Logic) CreateLink(file domain.File) (domain.File, error) {
	// check file existence and create a symlink.
	fi, err := lgc.Fsys.GetFileLstat(file.Link)
	if err != nil {
		// if file isn't found, create a symlink
		dir := filepath.Dir(file.Link)
		lgc.Printer.Printf("create directory %s\n", dir)
		if err := lgc.Fsys.MkdirAll(dir); err != nil {
			return file, err
		}
		lgc.Printer.Printf("create link %s -> %s\n", file.Link, file.Bin)
		if err := lgc.Fsys.MkLink(file.Bin, file.Link); err != nil {
			return file, err
		}
		file.Result.LinkCreated = true
		return file, nil
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		// if file is a directory, raise error
		return file, fmt.Errorf("%s has already existed and is a directory", file.Link)
	case mode&os.ModeNamedPipe != 0:
		// if file is a pipe, raise error
		return file, fmt.Errorf("%s has already existed and is a named pipe", file.Link)
	case mode.IsRegular():
		// if file is a regular file, remove it and create a symlink.
		return lgc.Logic.RemoveFileAndCreateLink(file)
	case mode&os.ModeSymlink != 0:
		return lgc.Logic.RecreateLink(file)
	default:
		return file, fmt.Errorf("unexpected file mode %s: %s", file.Link, mode.String())
	}
}

func (lgc *Logic) RemoveFileAndCreateLink(file domain.File) (domain.File, error) {
	lgc.Printer.Printf("remove %s\n", file.Link)
	if err := lgc.Fsys.RemoveFile(file.Link); err != nil {
		return file, err
	}
	file.Result.FileRemoved = true
	lgc.Printer.Printf("create link %s -> %s\n", file.Link, file.Bin)
	if err := lgc.Fsys.MkLink(file.Bin, file.Link); err != nil {
		return file, err
	}
	file.Result.Migrated = true
	return file, nil
}

func (lgc *Logic) RecreateLink(file domain.File) (domain.File, error) {
	// if file is a symlink but a dest is different, recreate a symlink.
	lnDest, err := lgc.Fsys.ReadLink(file.Link)
	if err != nil {
		return file, err
	}
	if file.Bin == lnDest {
		return file, nil
	}
	lgc.Printer.Printf("remove link %s -> %s\n", file.Link, lnDest)
	if err := lgc.Fsys.RemoveLink(file.Link); err != nil {
		return file, err
	}
	file.Result.LinkRemoved = true
	lgc.Printer.Printf("create link %s -> %s\n", file.Link, file.Bin)
	if err := lgc.Fsys.MkLink(file.Bin, file.Link); err != nil {
		return file, err
	}
	file.Result.Migrated = true
	return file, nil
}
