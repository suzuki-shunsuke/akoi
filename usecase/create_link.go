package usecase

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/suzuki-shunsuke/akoi/domain"
)

func createLink(pkg *domain.Package, file *domain.File, params *domain.InstallParams, methods *domain.InstallMethods) (*domain.FileResult, error) {
	fileResult := &domain.FileResult{}
	linkRelPath, err := filepath.Rel(filepath.Dir(file.Link), file.Bin)
	if err != nil {
		if params.Format != keyWordAnsible {
			fmt.Fprintln(os.Stderr, err)
		}
		return fileResult, err
	}
	if fi, err := methods.GetFileLstat(file.Link); err == nil {
		switch mode := fi.Mode(); {
		case mode.IsDir():
			if params.Format != keyWordAnsible {
				fmt.Fprintf(os.Stderr, "%s has already existed and is a directory\n", file.Link)
			}
			return fileResult, fmt.Errorf("%s has already existed and is a directory", file.Link)
		case mode&os.ModeNamedPipe != 0:
			if params.Format != keyWordAnsible {
				fmt.Fprintf(os.Stderr, "%s has already existed and is a named pipe\n", file.Link)
			}
			return fileResult, fmt.Errorf("%s has already existed and is a named pipe", file.Link)
		case mode.IsRegular():
			if params.Format != keyWordAnsible {
				fmt.Printf("remove %s\n", file.Link)
			}
			if err := methods.RemoveFile(file.Link); err != nil {
				if params.Format != keyWordAnsible {
					fmt.Fprintln(os.Stderr, err)
				}
				return fileResult, err
			}
			fileResult.Changed = true
			if params.Format != keyWordAnsible {
				fmt.Printf("create link %s -> %s\n", file.Link, linkRelPath)
			}
			if err := methods.MkLink(linkRelPath, file.Link); err != nil {
				if params.Format != keyWordAnsible {
					fmt.Fprintln(os.Stderr, err)
				}
				return fileResult, err
			}
			fileResult.Changed = true
			return fileResult, nil
		case mode&os.ModeSymlink != 0:
			lnDest, err := methods.ReadLink(file.Link)
			if err != nil {
				if params.Format != keyWordAnsible {
					fmt.Fprintln(os.Stderr, err)
				}
				return fileResult, err
			}
			if linkRelPath == lnDest {
				return fileResult, nil
			}
			if params.Format != keyWordAnsible {
				fmt.Printf("remove link %s -> %s\n", file.Link, lnDest)
			}
			if err := methods.RemoveLink(file.Link); err != nil {
				if params.Format != keyWordAnsible {
					fmt.Fprintln(os.Stderr, err)
				}
				return fileResult, err
			}
			fileResult.Changed = true
			if params.Format != keyWordAnsible {
				fmt.Printf("create link %s -> %s\n", file.Link, linkRelPath)
			}
			if err := methods.MkLink(linkRelPath, file.Link); err != nil {
				if params.Format != keyWordAnsible {
					fmt.Fprintln(os.Stderr, err)
				}
				return fileResult, err
			}
			return fileResult, nil
		default:
			return fileResult, fmt.Errorf("unexpected file mode %s: %s", file.Link, mode.String())
		}
	}
	if params.Format != keyWordAnsible {
		fmt.Printf("create link %s -> %s\n", file.Link, linkRelPath)
	}
	if err := methods.MkLink(linkRelPath, file.Link); err != nil {
		if params.Format != keyWordAnsible {
			fmt.Fprintln(os.Stderr, err)
		}
		return fileResult, err
	}
	fileResult.Changed = true
	return fileResult, nil
}
