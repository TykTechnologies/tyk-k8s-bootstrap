package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"tyk/tyk/bootstrap/data"
	"tyk/tyk/bootstrap/helpers"
	"tyk/tyk/bootstrap/readiness"
)

func main() {
	err := data.InitAppDataPostInstall()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = readiness.CheckIfRequiredDeploymentsAreReady()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tp := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: data.AppConfig.DashboardInsecureSkipVerify},
	}
	client := http.Client{Transport: tp}

	fmt.Println("Started creating dashboard org")
	err = helpers.CheckForExistingOrganisation(client)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Finished creating dashboard org")

	fmt.Println("Generating dashboard credentials")
	err = helpers.GenerateDashboardCredentials(client)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Finished generating dashboard credentials")

	fmt.Println("Started bootstrapping operator secret")
	if data.AppConfig.OperatorSecretEnabled {
		err = helpers.BootstrapTykOperatorSecret()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	fmt.Println("Finished bootstrapping operator secret")

	fmt.Println("Started bootstrapping portal secret")
	if data.AppConfig.DeveloperPortalSecretEnabled {
		err = helpers.BootstrapTykPortalSecret()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	fmt.Println("Started bootstrapping portal with requests to dashboard")
	if data.AppConfig.BootstrapPortal {
		err = helpers.BoostrapPortal(client)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	fmt.Println("Finished bootstrapping portal")

}
