package usecase

import (
	"testing"
)

func TestNewLogic(t *testing.T) {
	lgc := NewLogic()
	if lgc == nil {
		t.Fatal("logic is nil")
	}
}
