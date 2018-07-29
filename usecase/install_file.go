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

	// Check file
	if fi, err := methods.GetFileStat(dst); err == nil {
		if fi.Mode() == mode {
			return fileResult, nil
		}
		if params.Format != keyWordAnsible {
			fmt.Printf("chmod %s %s\n", mode.String(), dst)
		}
		if err := methods.Chmod(dst, mode); err != nil {
			if params.Format != keyWordAnsible {
				fmt.Fprintln(os.Stderr, err)
			}
			return fileResult, err
		}
		fileResult.Changed = true
		return fileResult, err
	}

	// Download
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

	// Create temporary directory
	tmpDir, err := methods.TempDir()
	if err != nil {
		if params.Format != keyWordAnsible {
			fmt.Fprintln(os.Stderr, err)
		}
		return fileResult, err
	}
	defer methods.RemoveAll(tmpDir)

	if pkg.ArchiveType != "unarchived" {
		arc := pkg.Archiver
		if arc == nil {
			if params.Format != keyWordAnsible {
				t := ustr
				if pkg.ArchiveType != "" {
					t = pkg.ArchiveType
				}
				fmt.Fprintf(os.Stderr, "failed to unarchive file: unsupported archive type: %s\n", t)
			}
			return fileResult, fmt.Errorf("failed to unarchive file: unsupported archive type")
		}
		// Unarchive
		if params.Format != keyWordAnsible {
			fmt.Printf("unarchive %s\n", pkg.Name)
		}
		if err := arc.Read(resp.Body, tmpDir); err != nil {
			if params.Format != keyWordAnsible {
				fmt.Fprintln(os.Stderr, err)
			}
			return fileResult, err
		}
	}

	for _, f := range pkg.Files {
		// Check file
		fi, err := methods.GetFileStat(f.Bin)
		if err != nil {
			// Create parent directory
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

			// Install
			if params.Format != keyWordAnsible {
				fmt.Printf("install %s\n", dst)
			}
			writer, err := methods.OpenFile(dst, os.O_RDWR|os.O_CREATE, mode)
			if err != nil {
				if params.Format != keyWordAnsible {
					fmt.Fprintf(os.Stderr, "failed to install %s: %s\n", dst, err)
				}
				return fileResult, err
			}
			defer writer.Close()
			if pkg.ArchiveType != "unarchived" {
				src, err := methods.Open(filepath.Join(tmpDir, f.Archive))
				if err != nil {
					if params.Format != keyWordAnsible {
						fmt.Fprintln(os.Stderr, err)
					}
					return fileResult, err
				}
				defer src.Close()
				if _, err := methods.Copy(writer, src); err != nil {
					if params.Format != keyWordAnsible {
						fmt.Fprintln(os.Stderr, err)
					}
					return fileResult, err
				}
			} else {
				if _, err := methods.Copy(writer, resp.Body); err != nil {
					if params.Format != keyWordAnsible {
						fmt.Fprintln(os.Stderr, err)
					}
					return fileResult, err
				}
			}
			fileResult.Changed = true
			continue
		}
		// Change file mode
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
	return fileResult, nil
}
