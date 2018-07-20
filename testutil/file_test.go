package testutil

import (
	"fmt"
	"testing"
)

func TestNewFakeExist(t *testing.T) {
	f := NewFakeExistFile(true)
	if !f("") {
		t.Fatal("result must be true")
	}
	f = NewFakeExistFile(false)
	if f("") {
		t.Fatal("result must be false")
	}
}

func TestNewFakeWrite(t *testing.T) {
	f := NewFakeWrite(nil)
	if err := f("", nil); err != nil {
		t.Fatal(err)
	}
	f = NewFakeWrite(fmt.Errorf("error"))
	if err := f("", nil); err == nil {
		t.Fatal("err must not be nil")
	}
}

func TestNewFakeMkdirAll(t *testing.T) {
	f := NewFakeMkdirAll(nil)
	if err := f(""); err != nil {
		t.Fatal(err)
	}
	f = NewFakeMkdirAll(fmt.Errorf("failed to create a directory"))
	if err := f(""); err == nil {
		t.Fatal("it must be failed to create a directory")
	}
}
