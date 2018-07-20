package domain

import (
	"testing"
)

func TestConfigGetBinPathTpl(t *testing.T) {
	cfg := &Config{
		LinkPath: "/usr/local/bin/{{.Name}}",
	}
	if _, err := cfg.GetLinkPathTpl(); err != nil {
		t.Fatal(err)
	}
	if _, err := cfg.GetLinkPathTpl(); err != nil {
		t.Fatal(err)
	}
}

func TestConfigGetLinkPathTpl(t *testing.T) {
	cfg := &Config{
		BinPath: "/usr/local/bin/{{.Name}}-{{.Version}}",
	}
	if _, err := cfg.GetBinPathTpl(); err != nil {
		t.Fatal(err)
	}
	if _, err := cfg.GetBinPathTpl(); err != nil {
		t.Fatal(err)
	}
}

func TestPackageGetArchiver(t *testing.T) {
	pkg := &Package{
		URL: "http://example.com/foo.zip",
	}
	arc, err := pkg.GetArchiver()
	if err != nil {
		t.Fatal(err)
	}
	if !arc.Match("zzz.zip") {
		t.Fatal("zzz.zip must be matched")
	}
}
