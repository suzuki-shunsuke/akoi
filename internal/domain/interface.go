package domain

import (
	"io"
	"net/http"
	"os"
)

type (
	// Chmod is the interface of os.Chmod .
	Chmod func(name string, mode os.FileMode) error
	// Copy is the interface of io.Copy .
	Copy func(dst io.Writer, src io.Reader) (int64, error)
	// Download downloads a file.
	Download func(url string) (*http.Response, error)
	// ExistFile is an interface to check file existence.
	ExistFile func(string) bool
	// ExpandEnv is an interface of os.ExpandEnv .
	ExpandEnv func(string) string
	// Fprintf is an interface of fmt.Fprintf .
	Fprintf func(w io.Writer, format string, a ...interface{}) (n int, err error)
	// Fprintln is an interface of fmt.Fprintln .
	Fprintln func(w io.Writer, a ...interface{}) (n int, err error)
	// GetArchiver returns an archiver for a given file
	GetArchiver func(fpath, ftype string) Archiver
	// GetFileStat returns a FileInfo.
	GetFileStat func(string) (os.FileInfo, error)
	// MkdirAll is an interface to create directories.
	MkdirAll func(string) error
	// MkLink creates a symbolic link.
	MkLink func(src, dst string) error
	// NewGzipReader creates a reader for gzip.
	NewGzipReader func(io.Reader) (io.ReadCloser, error)
	// Open opens a file.
	Open func(name string) (*os.File, error)
	// OpenFile opens a file.
	OpenFile func(name string, flag int, perm os.FileMode) (*os.File, error)
	// Printf is an interface of fmt.Printf .
	Printf func(format string, a ...interface{}) (n int, err error)
	// Println is an interface of fmt.Println .
	Println func(a ...interface{}) (n int, err error)
	// ReadConfigFile reads a configuration file.
	ReadConfigFile func(string) (*Config, error)
	// ReadLink gets a symbolic's destination path.
	ReadLink func(string) (string, error)
	// RemoveFile is an interface of os.Remove .
	RemoveFile func(string) error
	// TempDir creates a temporary directory.
	TempDir func() (string, error)
	// WriteFile is an interface to write test to file.
	WriteFile func(dest string, data []byte) error

	// Archiver is an interface to read an archive file.
	Archiver interface {
		Read(input io.Reader, destination string) error
	}
)
