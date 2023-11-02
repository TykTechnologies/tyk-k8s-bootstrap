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
	err := data.InitPostInstall()
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
		TLSClientConfig: &tls.Config{InsecureSkipVerify: data.BootstrapConf.InsecureSkipVerify},
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
	if data.BootstrapConf.OperatorKubernetesSecretName != "" {
		err = helpers.BootstrapTykOperatorSecret()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	fmt.Println("Finished bootstrapping operator secret")

	fmt.Println("Started bootstrapping portal secret")
	if data.BootstrapConf.DevPortalKubernetesSecretName != "" {
		err = helpers.BootstrapTykPortalSecret()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	fmt.Println("Started bootstrapping portal with requests to dashboard")
	if data.BootstrapConf.BootstrapPortal {
		err = helpers.BoostrapPortal(client)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	fmt.Println("Finished bootstrapping portal")

}
