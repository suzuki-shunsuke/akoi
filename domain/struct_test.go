package domain

import (
	"testing"
)

func TestPackageGetURL(t *testing.T) {
	exp := "foo"
	pkg := &Package{url: exp}
	act := pkg.GetURL()
	if act != exp {
		t.Fatalf(`pkg.GetURL() = "%s", wanted "%s"`, act, exp)
	}
}

func TestConfigSetup(t *testing.T) {
	cfg := &Config{
		BinPath:  "/usr/local/bin/{{.Name}}-{{.Version}}",
		LinkPath: "/usr/local/bin/{{.Name}}",
		Packages: map[string]Package{
			"consul": {
				URL:     "https://releases.hashicorp.com/consul/{{.Version}}/consul_{{.Version}}_darwin_amd64.zip",
				Version: "1.2.1",
				Files: []File{
					{
						Name:    "consul",
						Archive: "consul",
					},
				},
			},
		},
	}
	if err := cfg.Setup(); err != nil {
		t.Fatal(err)
	}
}
