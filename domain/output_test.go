package domain

import (
	"testing"
)

func TestResultString(t *testing.T) {
	exp := ""
	result := &Result{Msg: "foo"}
	params := &InstallParams{}
	act := result.String(params)
	if act != exp {
		t.Fatalf(`result.String(params) = "%s", wanted "%s"`, act, exp)
	}
	params.Format = "ansible"
	if result.String(params) == "" {
		t.Fatal("result.String(params) is empty")
	}
}
