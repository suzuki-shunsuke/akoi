package usecase

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
)

func getInstalledFiles(
	files []domain.File, methods domain.InstallMethods,
) []domain.File {
	installedFiles := []domain.File{}
	for _, file := range files {
		dst := file.Bin
		fileResult := file.Result
		mode := file.Mode
		if fi, err := methods.GetFileStat(dst); err == nil {
			if fi.Mode() == mode {
				continue
			}
			methods.Printf("chmod %s %s\n", mode.String(), dst)
			if err := methods.Chmod(dst, mode); err != nil {
				methods.Fprintln(os.Stderr, err)
				fileResult.Error = err.Error()
				continue
			}
			fileResult.Changed = true
			fileResult.ModeChanged = true
			continue
		}

		// Create parent directory
		dir := filepath.Dir(dst)
		if _, err := methods.GetFileStat(dir); err != nil {
			methods.Printf("create directory %s\n", dir)
			if err := methods.MkdirAll(dir); err != nil {
				methods.Fprintln(os.Stderr, err)
				fileResult.Error = err.Error()
				continue
			}
			fileResult.Changed = true
			fileResult.DirCreated = true
		}
		installedFiles = append(installedFiles, file)
	}
	return installedFiles
}

func installPackage(
	ctx context.Context, pkg *domain.Package, params domain.InstallParams,
	methods domain.InstallMethods,
) {
	installedFiles := getInstalledFiles(pkg.Files, methods)
	if len(installedFiles) != 0 {
		// Download
		ustr := pkg.URL.String()
		methods.Printf("downloading %s: %s\n", pkg.Name, ustr)
		c, cancel := context.WithCancel(ctx)
		defer cancel()
		resp, err := methods.Download(c, ustr)
		if err != nil {
			methods.Fprintln(os.Stderr, err)
			pkg.Result.Error = err.Error()
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			methods.Fprintf(os.Stderr, "failed to download %s from %s: %d\n", pkg.Name, ustr, resp.StatusCode)
			pkg.Result.Error = fmt.Sprintf(
				"failed to download %s from %s: %d", pkg.Name, ustr, resp.StatusCode)
			return
		}

		tmpDir := ""
		if pkg.Archived() {
			// Create temporary directory
			var err error
			tmpDir, err = methods.TempDir()
			if err != nil {
				methods.Fprintln(os.Stderr, err)
				pkg.Result.Error = err.Error()
				return
			}
			defer methods.RemoveAll(tmpDir)

			arc := pkg.Archiver
			if arc == nil {
				if params.Format != keyWordAnsible {
					t := ustr
					if pkg.ArchiveType != "" {
						t = pkg.ArchiveType
					}
					pkg.Result.Error = fmt.Sprintf("failed to unarchive file: unsupported archive type: %s\n", t)
					methods.Fprintf(os.Stderr, "failed to unarchive file: unsupported archive type: %s\n", t)
				}
				return
			}
			// Unarchive
			methods.Printf("unarchive %s\n", pkg.Name)
			if err := arc.Read(resp.Body, tmpDir); err != nil {
				pkg.Result.Error = err.Error()
				methods.Fprintln(os.Stderr, err)
				return
			}
		}

		for _, file := range installedFiles {
			// Install
			mode := file.Mode
			dst := file.Bin
			methods.Printf("install %s\n", dst)
			writer, err := methods.OpenFile(dst, os.O_RDWR|os.O_CREATE, mode)
			fileResult := file.Result
			if err != nil {
				methods.Fprintf(os.Stderr, "failed to install %s: %s\n", dst, err)
				fileResult.Error = err.Error()
				continue
			}
			defer writer.Close()
			if pkg.Archived() {
				src, err := methods.Open(filepath.Join(tmpDir, file.Archive))
				if err != nil {
					methods.Fprintln(os.Stderr, err)
					fileResult.Error = err.Error()
					continue
				}
				defer src.Close()
				if _, err := methods.Copy(writer, src); err != nil {
					methods.Fprintln(os.Stderr, err)
					fileResult.Error = err.Error()
					continue
				}
				fileResult.Installed = true
			} else {
				if pkg.ArchiveType == "Gzip" {
					reader, err := methods.NewGzipReader(resp.Body)
					if err != nil {
						methods.Fprintln(os.Stderr, err)
						fileResult.Error = err.Error()
						continue
					}
					defer reader.Close()
					if _, err := methods.Copy(writer, reader); err != nil {
						methods.Fprintln(os.Stderr, err)
						fileResult.Error = err.Error()
						continue
					}
				}
				if _, err := methods.Copy(writer, resp.Body); err != nil {
					methods.Fprintln(os.Stderr, err)
					fileResult.Error = err.Error()
					continue
				}
				fileResult.Installed = true
			}
		}
	}
	for i, file := range pkg.Files {
		if file.Result.Error != "" {
			continue
		}
		f, err := createLink(file, methods)
		if err != nil {
			if f.Result.Error == "" {
				f.Result.Error = err.Error()
			}
		}
		pkg.Files[i] = f
	}
}
