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
		exit(err)
	}

	licenseIsValid, err := pkg.ValidateDashboardLicense(conf.Tyk.DashboardLicense)
	if err != nil {
		exit(err)
	}

	if !licenseIsValid {
		exit(fmt.Errorf("provided license is invalid"))
	}

	fmt.Println("Pre-Hook bootstrapping succeeded, the provided license is valid!")
}

func exit(err error) {
	if err != nil {
		fmt.Printf("[ERROR]: %v", err)
		os.Exit(1)
	}
}
