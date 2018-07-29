package usecase

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/suzuki-shunsuke/akoi/domain"
)

func installFile(pkg *domain.Package, file *domain.File, params *domain.InstallParams, methods *domain.InstallMethods) (*domain.FileResult, error) {
	dst := file.Bin
	fileResult := &domain.FileResult{
		Name: file.Name,
	}
	mode := file.Mode
	if mode == 0 {
		mode = 0755
	}
	if fi, err := methods.GetFileStat(dst); err == nil {
		if fi.Mode() == mode {
			return fileResult, nil
		}
		if params.Format != keyWordAnsible {
			fmt.Printf("chmod %s %s\n", mode.String(), dst)
		}
		err := methods.Chmod(dst, mode)
		if err != nil && params.Format != keyWordAnsible {
			fmt.Fprintln(os.Stderr, err)
		}
		return fileResult, err
	}

	ustr := pkg.URL.String()
	if params.Format != keyWordAnsible {
		fmt.Printf("downloading %s: %s\n", pkg.Name, ustr)
	}
	resp, err := methods.Download(ustr)
	if err != nil {
		if params.Format != keyWordAnsible {
			fmt.Fprintln(os.Stderr, err)
		}
		return fileResult, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		if params.Format != keyWordAnsible {
			fmt.Fprintf(os.Stderr, "failed to download %s from %s: %d\n", pkg.Name, ustr, resp.StatusCode)
		}
		return fileResult, fmt.Errorf(
			"failed to download %s from %s: %d", pkg.Name, ustr, resp.StatusCode)
	}
	tmpDir, err := methods.TempDir()
	if err != nil {
		if params.Format != keyWordAnsible {
			fmt.Fprintln(os.Stderr, err)
		}
		return fileResult, err
	}
	defer methods.RemoveAll(tmpDir)
	arc := pkg.Archiver
	// TODO support not archived file
	if arc != nil {
		if params.Format != keyWordAnsible {
			fmt.Printf("unarchive %s\n", pkg.Name)
		}
		if err := arc.Read(resp.Body, tmpDir); err != nil {
			if params.Format != keyWordAnsible {
				fmt.Fprintln(os.Stderr, err)
			}
			return fileResult, err
		}
		for _, f := range pkg.Files {
			fi, err := methods.GetFileStat(f.Bin)
			if err != nil {
				dir := filepath.Dir(f.Bin)
				if _, err := methods.GetFileStat(dir); err != nil {
					if params.Format != keyWordAnsible {
						fmt.Printf("create directory %s\n", dir)
					}
					if err := methods.MkdirAll(dir); err != nil {
						if params.Format != keyWordAnsible {
							fmt.Fprintln(os.Stderr, err)
						}
						return fileResult, err
					}
					fileResult.Changed = true
				}
				if params.Format != keyWordAnsible {
					fmt.Printf("install %s\n", dst)
				}
				if err := methods.CopyFile(filepath.Join(tmpDir, f.Archive), dst); err != nil {
					if params.Format != keyWordAnsible {
						fmt.Fprintln(os.Stderr, err)
					}
					return fileResult, err
				}
				fileResult.Changed = true
			}
			if err == nil && fi.Mode() == mode {
				continue
			}
			if err := methods.Chmod(dst, mode); err != nil {
				if params.Format != keyWordAnsible {
					fmt.Fprintln(os.Stderr, err)
				}
				return fileResult, err
			}
			fileResult.Changed = true
		}
	}
	return fileResult, nil
}
