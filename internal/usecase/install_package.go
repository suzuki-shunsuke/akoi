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
) domain.Package {
	installedFiles := lgc.Logic.GetInstalledFiles(pkg.Files)
	if len(installedFiles) != 0 {
		// Download
		ustr := pkg.URL.String()
		lgc.Printer.Printf("downloading %s: %s\n", pkg.Name, ustr)
		body, err := lgc.Downloader.Download(ctx, ustr, pkg.NumOfDLPartitions)
		if err != nil {
			lgc.Printer.Fprintln(os.Stderr, err)
			pkg.Result.Error = err.Error()
			return pkg
		}
		defer body.Close()
		tmpDir := ""
		if pkg.Archived() {
			// Create temporary directory
			var err error
			tmpDir, err = lgc.Fsys.TempDir()
			if err != nil {
				lgc.Printer.Fprintln(os.Stderr, err)
				pkg.Result.Error = err.Error()
				return pkg
			}
			defer lgc.Fsys.RemoveAll(tmpDir)

			arc := pkg.Archiver
			if arc == nil {
				if params.Format != keyWordAnsible {
					t := ustr
					if pkg.ArchiveType != "" {
						t = pkg.ArchiveType
					}
					pkg.Result.Error = fmt.Sprintf("failed to unarchive file: unsupported archive type: %s\n", t)
					lgc.Printer.Fprintf(os.Stderr, "failed to unarchive file: unsupported archive type: %s\n", t)
				}
				return pkg
			}
			// Unarchive
			lgc.Printer.Printf("unarchive %s\n", pkg.Name)
			if err := arc.Read(body, tmpDir); err != nil {
				pkg.Result.Error = err.Error()
				lgc.Printer.Fprintln(os.Stderr, err)
				return pkg
			}
		}

		for _, file := range installedFiles {
			// Install
			mode := file.Mode
			dst := file.Bin
			lgc.Printer.Printf("install %s\n", dst)
			writer, err := lgc.Fsys.OpenFile(dst, os.O_RDWR|os.O_CREATE, mode)
			if err != nil {
				lgc.Printer.Fprintf(os.Stderr, "failed to install %s: %s\n", dst, err)
				file.Result.Error = err.Error()
				continue
			}
			defer writer.Close()
			if pkg.Archived() {
				src, err := lgc.Fsys.Open(filepath.Join(tmpDir, file.Archive))
				if err != nil {
					lgc.Printer.Fprintln(os.Stderr, err)
					file.Result.Error = err.Error()
					continue
				}
				defer src.Close()
				if _, err := lgc.Fsys.Copy(writer, src); err != nil {
					lgc.Printer.Fprintln(os.Stderr, err)
					file.Result.Error = err.Error()
					continue
				}
				file.Result.Installed = true
			} else {
				if pkg.ArchiveType == "Gzip" {
					reader, err := lgc.GetGzipReader.Get(body)
					if err != nil {
						lgc.Printer.Fprintln(os.Stderr, err)
						file.Result.Error = err.Error()
						continue
					}
					defer reader.Close()
					if _, err := lgc.Fsys.Copy(writer, reader); err != nil {
						lgc.Printer.Fprintln(os.Stderr, err)
						file.Result.Error = err.Error()
						continue
					}
				}
				if _, err := lgc.Fsys.Copy(writer, body); err != nil {
					lgc.Printer.Fprintln(os.Stderr, err)
					file.Result.Error = err.Error()
					continue
				}
				file.Result.Installed = true
			}
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
	return pkg
}
