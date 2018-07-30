package registry

import (
	"testing"

	"github.com/suzuki-shunsuke/akoi/domain"
)

func TestNewInitMethods(t *testing.T) {
	methods := NewInitMethods()
	if methods == nil {
		t.Fatal("methods is nil")
	}
}

func TestNewInstallMethods(t *testing.T) {
	methods := NewInstallMethods(&domain.InstallParams{DryRun: true})
	if methods == nil {
		t.Fatal("methods is nil")
	}
	methods = NewInstallMethods(&domain.InstallParams{DryRun: false})
	if methods == nil {
		t.Fatal("methods is nil")
	}
}
