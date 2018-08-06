package usecase

import (
	"os"

	"github.com/suzuki-shunsuke/akoi/domain"
	"github.com/suzuki-shunsuke/akoi/util"
)

// List intalls binraries.
func List(params *domain.ListParams, methods *domain.ListMethods) {
	if err := util.ValidateStruct(methods); err != nil {
		if methods.Fprintln != nil {
			methods.Fprintln(os.Stderr, err)
		}
	}
	cfg, err := methods.ReadConfigFile(params.ConfigFilePath)
	if err != nil {
		methods.Fprintln(os.Stderr, err)
	}
	if err := setupConfig(cfg, &domain.SetupConfigMethods{
		GetArchiver: methods.GetArchiver,
	}); err != nil {
		methods.Fprintln(os.Stderr, err)
	}
	numOfPkgs := len(cfg.Packages)
	if numOfPkgs == 0 {
		return
	}
	for _, pkg := range cfg.Packages {
		listPackage(&pkg, params, methods)
	}
}
