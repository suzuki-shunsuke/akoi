package usecase

import (
	"fmt"
	"os"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
)

func createLink(
	file *domain.File, methods *domain.InstallMethods,
) error {
	// check file existence and create a symlink.
	fileResult := file.Result
	fi, err := methods.GetFileLstat(file.Link)
	if err != nil {
		// if file isn't found, create a symlink
		methods.Printf("create link %s -> %s\n", file.Link, file.Bin)
		if err := methods.MkLink(file.Bin, file.Link); err != nil {
			methods.Fprintln(os.Stderr, err)
			return err
		}
		fileResult.Changed = true
		return nil
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		// if file is a directory, raise error
		methods.Fprintf(
			os.Stderr, "%s has already existed and is a directory\n", file.Link)
		return fmt.Errorf("%s has already existed and is a directory", file.Link)
	case mode&os.ModeNamedPipe != 0:
		// if file is a pipe, raise error
		methods.Fprintf(
			os.Stderr, "%s has already existed and is a named pipe\n", file.Link)
		return fmt.Errorf("%s has already existed and is a named pipe", file.Link)
	case mode.IsRegular():
		// if file is a regular file, remove it and create a symlink.
		methods.Printf("remove %s\n", file.Link)
		if err := methods.RemoveFile(file.Link); err != nil {
			methods.Fprintln(os.Stderr, err)
			return err
		}
		fileResult.Changed = true
		fileResult.FileRemoved = true
		methods.Printf("create link %s -> %s\n", file.Link, file.Bin)
		if err := methods.MkLink(file.Bin, file.Link); err != nil {
			methods.Fprintln(os.Stderr, err)
			return err
		}
		fileResult.Migrated = true
		fileResult.Changed = true
		return nil
	case mode&os.ModeSymlink != 0:
		// if file is a symlink but a dest is different, recreate a symlink.
		lnDest, err := methods.ReadLink(file.Link)
		if err != nil {
			methods.Fprintln(os.Stderr, err)
			return err
		}
		if file.Bin == lnDest {
			return nil
		}
		methods.Printf("remove link %s -> %s\n", file.Link, lnDest)
		if err := methods.RemoveLink(file.Link); err != nil {
			methods.Fprintln(os.Stderr, err)
			return err
		}
		fileResult.Changed = true
		methods.Printf("create link %s -> %s\n", file.Link, file.Bin)
		if err := methods.MkLink(file.Bin, file.Link); err != nil {
			methods.Fprintln(os.Stderr, err)
			return err
		}
		fileResult.Migrated = true
		return nil
	default:
		return fmt.Errorf("unexpected file mode %s: %s", file.Link, mode.String())
	}
}
