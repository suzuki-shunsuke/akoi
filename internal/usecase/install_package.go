package usecase

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
)

func getInstalledFiles(
	files []domain.File, fsys domain.FileSystem, printer domain.Printer,
) []domain.File {
	installedFiles := []domain.File{}
	for _, file := range files {
		dst := file.Bin
		mode := file.Mode
		if fi, err := fsys.GetFileStat(dst); err == nil {
			if fi.Mode() == mode {
				continue
			}
			printer.Printf("chmod %s %s\n", mode.String(), dst)
			if err := fsys.Chmod(dst, mode); err != nil {
				printer.Fprintln(os.Stderr, err)
				file.Result.Error = err.Error()
				continue
			}
			file.Result.ModeChanged = true
			continue
		}

		// Create parent directory
		dir := filepath.Dir(dst)
		if _, err := fsys.GetFileStat(dir); err != nil {
			printer.Printf("create directory %s\n", dir)
			if err := fsys.MkdirAll(dir); err != nil {
				printer.Fprintln(os.Stderr, err)
				file.Result.Error = err.Error()
				continue
			}
			file.Result.DirCreated = true
		}
		installedFiles = append(installedFiles, file)
	}
	return installedFiles
}

func (lgc *logic) InstallPackage(
	ctx context.Context, pkg domain.Package, params domain.InstallParams,
	fsys domain.FileSystem, printer domain.Printer, downloader domain.Downloader, getGzipReader domain.GetGzipReader,
) domain.Package {
	installedFiles := getInstalledFiles(pkg.Files, fsys, printer)
	if len(installedFiles) != 0 {
		// Download
		ustr := pkg.URL.String()
		printer.Printf("downloading %s: %s\n", pkg.Name, ustr)
		body, err := downloader.Download(ctx, ustr, pkg.NumOfDLPartitions)
		if err != nil {
			printer.Fprintln(os.Stderr, err)
			pkg.Result.Error = err.Error()
			return pkg
		}
		defer body.Close()
		tmpDir := ""
		if pkg.Archived() {
			// Create temporary directory
			var err error
			tmpDir, err = fsys.TempDir()
			if err != nil {
				printer.Fprintln(os.Stderr, err)
				pkg.Result.Error = err.Error()
				return pkg
			}
			defer fsys.RemoveAll(tmpDir)

			arc := pkg.Archiver
			if arc == nil {
				if params.Format != keyWordAnsible {
					t := ustr
					if pkg.ArchiveType != "" {
						t = pkg.ArchiveType
					}
					pkg.Result.Error = fmt.Sprintf("failed to unarchive file: unsupported archive type: %s\n", t)
					printer.Fprintf(os.Stderr, "failed to unarchive file: unsupported archive type: %s\n", t)
				}
				return pkg
			}
			// Unarchive
			printer.Printf("unarchive %s\n", pkg.Name)
			if err := arc.Read(body, tmpDir); err != nil {
				pkg.Result.Error = err.Error()
				printer.Fprintln(os.Stderr, err)
				return pkg
			}
		}

		for _, file := range installedFiles {
			// Install
			mode := file.Mode
			dst := file.Bin
			printer.Printf("install %s\n", dst)
			writer, err := fsys.OpenFile(dst, os.O_RDWR|os.O_CREATE, mode)
			if err != nil {
				printer.Fprintf(os.Stderr, "failed to install %s: %s\n", dst, err)
				file.Result.Error = err.Error()
				continue
			}
			defer writer.Close()
			if pkg.Archived() {
				src, err := fsys.Open(filepath.Join(tmpDir, file.Archive))
				if err != nil {
					printer.Fprintln(os.Stderr, err)
					file.Result.Error = err.Error()
					continue
				}
				defer src.Close()
				if _, err := fsys.Copy(writer, src); err != nil {
					printer.Fprintln(os.Stderr, err)
					file.Result.Error = err.Error()
					continue
				}
				file.Result.Installed = true
			} else {
				if pkg.ArchiveType == "Gzip" {
					reader, err := getGzipReader.Get(body)
					if err != nil {
						printer.Fprintln(os.Stderr, err)
						file.Result.Error = err.Error()
						continue
					}
					defer reader.Close()
					if _, err := fsys.Copy(writer, reader); err != nil {
						printer.Fprintln(os.Stderr, err)
						file.Result.Error = err.Error()
						continue
					}
				}
				if _, err := fsys.Copy(writer, body); err != nil {
					printer.Fprintln(os.Stderr, err)
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
		f, err := createLink(file, fsys, printer)
		if err != nil {
			if f.Result.Error == "" {
				f.Result.Error = err.Error()
			}
		}
		pkg.Files[i] = f
	}
	return pkg
}
