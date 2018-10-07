package registry

import (
	"testing"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
)

func TestNewInitMethods(t *testing.T) {
	methods := NewInitMethods()
	if methods == nil {
		t.Fatal("methods is nil")
	}
}

func TestNewInstallMethods(t *testing.T) {
	NewInstallMethods(&domain.InstallParams{DryRun: true})
	NewInstallMethods(&domain.InstallParams{DryRun: false})
}
