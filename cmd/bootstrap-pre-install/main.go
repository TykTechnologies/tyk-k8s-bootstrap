package main

import (
	"fmt"
	"os"
	"tyk/tyk/bootstrap/pkg"
	"tyk/tyk/bootstrap/pkg/config"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		os.Exit(1)
	}

	licenseIsValid, err := pkg.ValidateDashboardLicense(conf.Tyk.DashboardLicense)
	if err != nil {
		os.Exit(1)
	}

	if !licenseIsValid {
		//return errors.New("provided license is invalid")
		os.Exit(1)
	}

	fmt.Println("Pre-Hook bootstrapping succeeded, the provided license is valid!")
}
