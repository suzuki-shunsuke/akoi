package infra

import (
	"runtime"
)

type (
	// Runtime implements domain.Runtime
	Runtime struct{}
)

// OS returns the running program's operating system target.
func (rt *Runtime) OS() string {
	return runtime.GOOS
}

// Arch returns the running program's architecture target.
func (rt *Runtime) Arch() string {
	return runtime.GOARCH
}

// NumCPU returns the number of logical CPUs usable by the current process.
func (rt *Runtime) NumCPU() int {
	return runtime.NumCPU()
}
