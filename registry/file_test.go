package registry

import (
	"testing"
)

func TestNewInitMethods(t *testing.T) {
	methods := NewInitMethods()
	if methods == nil {
		t.Fatal("methods is nil")
	}
}

func TestNewInstallMethods(t *testing.T) {
	methods := NewInstallMethods(false)
	if methods == nil {
		t.Fatal("methods is nil")
	}
	methods = NewInstallMethods(true)
	if methods == nil {
		t.Fatal("methods is nil")
	}
}
