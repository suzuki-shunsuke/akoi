package domain

import (
	"context"
	"io"
	"os"
)

type (
	// GetArchiver returns an archiver for a given file
	GetArchiver interface {
		Get(fpath, ftype string) Archiver
	}

	// GetGzipReader creates a reader for gzip.
	GetGzipReader interface {
		Get(io.Reader) (io.ReadCloser, error)
	}

	// Archiver is an interface to read an archive file.
	Archiver interface {
		Read(input io.Reader, destination string) error
	}

	// FileSystem abstracts file system's operation.
	FileSystem interface {
		Chmod(name string, mode os.FileMode) error
		Copy(dst io.Writer, src io.Reader) (int64, error)
		ExistFile(string) bool
		ExpandEnv(string) string
		GetFileLstat(string) (os.FileInfo, error)
		GetFileStat(string) (os.FileInfo, error)
		MkdirAll(string) error
		MkLink(src, dst string) error
		Open(name string) (*os.File, error)
		OpenFile(name string, flag int, perm os.FileMode) (*os.File, error)
		ReadLink(string) (string, error)
		RemoveAll(string) error
		RemoveFile(string) error
		RemoveLink(string) error
		TempDir() (string, error)
		WriteFile(dest string, data []byte) error
	}

	// Printer outputs messages.
	Printer interface {
		Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error)
		Fprintln(w io.Writer, a ...interface{}) (n int, err error)
		Printf(format string, a ...interface{}) (n int, err error)
		Println(a ...interface{}) (n int, err error)
	}

	// ConfigReader reads the configuration file.
	ConfigReader interface {
		Read(string) (Config, error)
	}

	// Downloader downloads a file.
	Downloader interface {
		Download(ctx context.Context, uri string, numOfDLPartitions int) (io.ReadCloser, error)
	}

	// Logic abstracts application logic.
	Logic interface {
		Install(ctx context.Context, params InstallParams, printer Printer, cfgReader ConfigReader, getArchiver GetArchiver, downloader Downloader, getGzipReader GetGzipReader) Result
		InstallPackage(ctx context.Context, pkg Package, params InstallParams, fsys FileSystem, printer Printer, downloader Downloader, getGzipReader GetGzipReader) Package
		GetInstalledFiles(files []File, fsys FileSystem, printer Printer) []File
		CreateLink(file File, fsys FileSystem, printer Printer) (File, error)
		SetupConfig(cfg Config, fsys FileSystem, getArchiver GetArchiver) (Config, error)
	}
)
