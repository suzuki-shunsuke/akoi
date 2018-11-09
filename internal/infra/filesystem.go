package infra

import (
	"io"
	"io/ioutil"
	"os"
)

type (
	// FileSystem implements domain.FileSystem .
	FileSystem struct{}
)

// Chmod implements domain.FileSystem .
func (fsys FileSystem) Chmod(name string, mode os.FileMode) error {
	return os.Chmod(name, mode)
}

// Copy implements domain.FileSystem .
func (fsys FileSystem) Copy(dst io.Writer, src io.Reader) (int64, error) {
	return io.Copy(dst, src)
}

// ExistFile implements domain.FileSystem .
func (fsys FileSystem) ExistFile(dst string) bool {
	_, err := os.Stat(dst)
	return err == nil
}

// ExpandEnv implements domain.FileSystem .
func (fsys FileSystem) ExpandEnv(p string) string {
	return os.ExpandEnv(p)
}

// GetFileLstat implements domain.FileSystem .
func (fsys FileSystem) GetFileLstat(p string) (os.FileInfo, error) {
	return os.Lstat(p)
}

// GetFileStat implements domain.FileSystem .
func (fsys FileSystem) GetFileStat(p string) (os.FileInfo, error) {
	return os.Stat(p)
}

// MkdirAll implements domain.FileSystem .
func (fsys FileSystem) MkdirAll(dst string) error {
	return os.MkdirAll(dst, 0775)
}

// MkLink implements domain.FileSystem .
func (fsys FileSystem) MkLink(src, dst string) error {
	return os.Symlink(src, dst)
}

// Open implements domain.FileSystem .
func (fsys FileSystem) Open(name string) (*os.File, error) {
	return os.Open(name)
}

// OpenFile implements domain.FileSystem .
func (fsys FileSystem) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

// ReadLink implements domain.FileSystem .
func (fsys FileSystem) ReadLink(p string) (string, error) {
	return os.Readlink(p)
}

// RemoveAll implements domain.FileSystem .
func (fsys FileSystem) RemoveAll(p string) error {
	return os.RemoveAll(p)
}

// RemoveFile implements domain.FileSystem .
func (fsys FileSystem) RemoveFile(p string) error {
	return os.Remove(p)
}

// RemoveLink implements domain.FileSystem .
func (fsys FileSystem) RemoveLink(p string) error {
	return os.Remove(p)
}

// TempDir implements domain.FileSystem .
func (fsys FileSystem) TempDir() (string, error) {
	return ioutil.TempDir("", "")
}

// WriteFile implements domain.FileSystem .
func (fsys FileSystem) WriteFile(dst string, data []byte) error {
	return ioutil.WriteFile(dst, data, 0644)
}
