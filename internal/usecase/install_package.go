package usecase

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
)

func (lgc *Logic) GetInstalledFiles(files []domain.File) []domain.File {
	installedFiles := []domain.File{}
	for _, file := range files {
		dst := file.Bin
		mode := file.Mode
		if fi, err := lgc.Fsys.GetFileStat(dst); err == nil {
			if fi.Mode() == mode {
				continue
			}
			lgc.Printer.Printf("chmod %s %s\n", mode.String(), dst)
			if err := lgc.Fsys.Chmod(dst, mode); err != nil {
				lgc.Printer.Fprintln(os.Stderr, err)
				file.Result.Error = err.Error()
				continue
			}
			file.Result.ModeChanged = true
			continue
		}

		// Create parent directory
		dir := filepath.Dir(dst)
		if _, err := lgc.Fsys.GetFileStat(dir); err != nil {
			lgc.Printer.Printf("create directory %s\n", dir)
			if err := lgc.Fsys.MkdirAll(dir); err != nil {
				lgc.Printer.Fprintln(os.Stderr, err)
				file.Result.Error = err.Error()
				continue
			}
			file.Result.DirCreated = true
		}
		installedFiles = append(installedFiles, file)
	}
	return installedFiles
}

func (lgc *Logic) InstallPackage(
	ctx context.Context, pkg domain.Package, params domain.InstallParams,
) (domain.Package, error) {
	installedFiles := lgc.Logic.GetInstalledFiles(pkg.Files)
	if len(installedFiles) != 0 {
		// Download
		ustr := pkg.URL.String()
		lgc.Printer.Printf("downloading %s: %s\n", pkg.Name, ustr)
		body, err := lgc.Downloader.Download(ctx, ustr, pkg.NumOfDLPartitions)
		if err != nil {
			lgc.Printer.Fprintln(os.Stderr, err)
			return pkg, err
		}
		defer body.Close()
		tmpDir := ""
		if pkg.Archived() {
			// Create temporary directory
			var err error
			tmpDir, err = lgc.Fsys.TempDir()
			if err != nil {
				lgc.Printer.Fprintln(os.Stderr, err)
				return pkg, err
			}
			defer lgc.Fsys.RemoveAll(tmpDir)

			arc := pkg.Archiver
			if arc == nil {
				t := ustr
				if pkg.ArchiveType != "" {
					t = pkg.ArchiveType
				}
				if params.Format != keyWordAnsible {
					lgc.Printer.Fprintf(os.Stderr, "failed to unarchive file: unsupported archive type: %s\n", t)
				}
				return pkg, fmt.Errorf("failed to unarchive file: unsupported archive type: %s", t)
			}
			// Unarchive
			lgc.Printer.Printf("unarchive %s\n", pkg.Name)
			if err := arc.Read(body, tmpDir); err != nil {
				lgc.Printer.Fprintln(os.Stderr, err)
				return pkg, err
			}
		}

		for _, file := range installedFiles {
			if err := lgc.Logic.InstallFile(&file, pkg, params, tmpDir, body); err != nil {
				lgc.Printer.Fprintf(os.Stderr, "failed to install %s: %s\n", file.Bin, err)
				file.Result.Error = err.Error()
			}
			file.Result.Installed = true
		}
	}
	for i, file := range pkg.Files {
		if file.Result.Error != "" {
			continue
		}
		f, err := lgc.Logic.CreateLink(file)
		if err != nil {
			if f.Result.Error == "" {
				f.Result.Error = err.Error()
			}
		}
		pkg.Files[i] = f
	}
	return pkg, nil
}
