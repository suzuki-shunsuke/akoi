package util

import (
	"testing"
	"text/template"
)

func TestRenderTpl(t *testing.T) {
	exp := "foo"
	tpl, err := template.New("link_path").Parse(exp)
	if err != nil {
		t.Fatal(err)
	}
	s, err := RenderTpl(tpl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if s != exp {
		t.Fatalf(`s = "%s", wanted "%s"`, s, exp)
	}
}
