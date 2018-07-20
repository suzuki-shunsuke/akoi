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
	methods := NewInstallMethods()
	if methods == nil {
		t.Fatal("methods is nil")
	}
}
