package main

import (
	"fmt"
	"net/http"
	"os"
	"tyk/tyk/bootstrap/data"
	"tyk/tyk/bootstrap/helpers"
	"tyk/tyk/bootstrap/license"
	"tyk/tyk/bootstrap/readiness"
)

func main() {
	err := data.InitAppDataPostInstall()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dashboardLicenseKey, err := license.GetDashboardLicense()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	licenseIsValid, err := license.ValidateDashboardLicense(dashboardLicenseKey)
	if err != nil {
		fmt.Println(err)
	}
	if licenseIsValid {
		fmt.Println("license is valid")
	} else {
		fmt.Println("license is invalid")
		os.Exit(1)
	}

	err = readiness.CheckIfDeploymentsAreReady()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client := http.Client{}

	err = helpers.CheckForExistingOrganisation(client)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = helpers.GenerateCredentials(client)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if data.AppConfig.BootstrapPortal {
		err = helpers.BoostrapPortal(client)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if data.AppConfig.OperatorSecretEnabled {
		err = helpers.BootstrapTykOperatorSecret()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if data.AppConfig.EnterprisePortalSecretEnabled {
		err = helpers.BootstrapTykEnterprisePortalSecret()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
