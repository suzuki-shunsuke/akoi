package domain

import (
	"testing"
)

func TestResultString(t *testing.T) {
	exp := ""
	result := &Result{Msg: "foo"}
	act := result.String("")
	if act != exp {
		t.Fatalf(`result.String(params) = "%s", wanted "%s"`, act, exp)
	}
	if result.String("ansible") == "" {
		t.Fatal("result.String(params) is empty")
	}
}
