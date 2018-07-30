package usecase

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/suzuki-shunsuke/akoi/domain"
)

func createLink(pkg *domain.Package, file *domain.File, params *domain.InstallParams, methods *domain.InstallMethods) error {
	fileResult := file.Result
	linkRelPath, err := filepath.Rel(filepath.Dir(file.Link), file.Bin)
	if err != nil {
		methods.Fprintln(os.Stderr, err)
		return err
	}
	if fi, err := methods.GetFileLstat(file.Link); err == nil {
		switch mode := fi.Mode(); {
		case mode.IsDir():
			methods.Fprintf(os.Stderr, "%s has already existed and is a directory\n", file.Link)
			return fmt.Errorf("%s has already existed and is a directory", file.Link)
		case mode&os.ModeNamedPipe != 0:
			methods.Fprintf(os.Stderr, "%s has already existed and is a named pipe\n", file.Link)
			return fmt.Errorf("%s has already existed and is a named pipe", file.Link)
		case mode.IsRegular():
			methods.Printf("remove %s\n", file.Link)
			if err := methods.RemoveFile(file.Link); err != nil {
				methods.Fprintln(os.Stderr, err)
				return err
			}
			fileResult.Changed = true
			fileResult.FileRemoved = true
			methods.Printf("create link %s -> %s\n", file.Link, linkRelPath)
			if err := methods.MkLink(linkRelPath, file.Link); err != nil {
				methods.Fprintln(os.Stderr, err)
				return err
			}
			fileResult.Migrated = true
			fileResult.Changed = true
			return nil
		case mode&os.ModeSymlink != 0:
			lnDest, err := methods.ReadLink(file.Link)
			if err != nil {
				methods.Fprintln(os.Stderr, err)
				return err
			}
			if linkRelPath == lnDest {
				return nil
			}
			methods.Printf("remove link %s -> %s\n", file.Link, lnDest)
			if err := methods.RemoveLink(file.Link); err != nil {
				methods.Fprintln(os.Stderr, err)
				return err
			}
			fileResult.Changed = true
			methods.Printf("create link %s -> %s\n", file.Link, linkRelPath)
			if err := methods.MkLink(linkRelPath, file.Link); err != nil {
				methods.Fprintln(os.Stderr, err)
				return err
			}
			fileResult.Migrated = true
			return nil
		default:
			return fmt.Errorf("unexpected file mode %s: %s", file.Link, mode.String())
		}
	}
	methods.Printf("create link %s -> %s\n", file.Link, linkRelPath)
	if err := methods.MkLink(linkRelPath, file.Link); err != nil {
		methods.Fprintln(os.Stderr, err)
		return err
	}
	fileResult.Changed = true
	return nil
}
