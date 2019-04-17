package usecase

import (
	"fmt"
	"os"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
)

func (lgc *logic) CreateLink(
	file domain.File,
	fsys domain.FileSystem, printer domain.Printer,
) (domain.File, error) {
	// check file existence and create a symlink.
	fi, err := fsys.GetFileLstat(file.Link)
	if err != nil {
		// if file isn't found, create a symlink
		printer.Printf("create link %s -> %s\n", file.Link, file.Bin)
		if err := fsys.MkLink(file.Bin, file.Link); err != nil {
			printer.Fprintln(os.Stderr, err)
			return file, err
		}
		file.Result.LinkCreated = true
		return file, nil
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		// if file is a directory, raise error
		printer.Fprintf(
			os.Stderr, "%s has already existed and is a directory\n", file.Link)
		return file, fmt.Errorf("%s has already existed and is a directory", file.Link)
	case mode&os.ModeNamedPipe != 0:
		// if file is a pipe, raise error
		printer.Fprintf(
			os.Stderr, "%s has already existed and is a named pipe\n", file.Link)
		return file, fmt.Errorf("%s has already existed and is a named pipe", file.Link)
	case mode.IsRegular():
		// if file is a regular file, remove it and create a symlink.
		printer.Printf("remove %s\n", file.Link)
		if err := fsys.RemoveFile(file.Link); err != nil {
			printer.Fprintln(os.Stderr, err)
			return file, err
		}
		file.Result.FileRemoved = true
		printer.Printf("create link %s -> %s\n", file.Link, file.Bin)
		if err := fsys.MkLink(file.Bin, file.Link); err != nil {
			printer.Fprintln(os.Stderr, err)
			return file, err
		}
		file.Result.Migrated = true
		return file, nil
	case mode&os.ModeSymlink != 0:
		// if file is a symlink but a dest is different, recreate a symlink.
		lnDest, err := fsys.ReadLink(file.Link)
		if err != nil {
			printer.Fprintln(os.Stderr, err)
			return file, err
		}
		if file.Bin == lnDest {
			return file, nil
		}
		printer.Printf("remove link %s -> %s\n", file.Link, lnDest)
		if err := fsys.RemoveLink(file.Link); err != nil {
			printer.Fprintln(os.Stderr, err)
			return file, err
		}
		file.Result.LinkRemoved = true
		printer.Printf("create link %s -> %s\n", file.Link, file.Bin)
		if err := fsys.MkLink(file.Bin, file.Link); err != nil {
			printer.Fprintln(os.Stderr, err)
			return file, err
		}
		file.Result.Migrated = true
		return file, nil
	default:
		return file, fmt.Errorf("unexpected file mode %s: %s", file.Link, mode.String())
	}
}
