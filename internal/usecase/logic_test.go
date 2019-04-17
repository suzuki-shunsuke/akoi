package usecase

import (
	"testing"
)

func TestNewLogic(t *testing.T) {
	lgc := NewLogic(nil)
	if lgc == nil {
		t.Fatal("logic is nil")
	}
}
