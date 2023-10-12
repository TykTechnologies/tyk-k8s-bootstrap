package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"tyk/tyk/bootstrap/k8s"
	"tyk/tyk/bootstrap/tyk"
	"tyk/tyk/bootstrap/tyk/data"
)

func main() {
	k8sClient, err := k8s.NewClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	appArgs, err := data.InitAppDataPostInstall()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	k8sClient.AppArgs = appArgs

	err = k8sClient.CheckIfRequiredDeploymentsAreReady()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tp := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: appArgs.DashboardInsecureSkipVerify},
	}
	client := http.Client{Transport: tp}

	tykSvc := tyk.NewTykService(client, appArgs)

	fmt.Println("Started creating dashboard org")
	err = tykSvc.CheckForExistingOrganisation()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Finished creating dashboard org")

	fmt.Println("Generating dashboard credentials")
	err = tykSvc.GenerateDashboardCredentials()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Finished generating dashboard credentials")

	fmt.Println("Started bootstrapping operator secret")
	if appArgs.OperatorSecretEnabled {
		err = k8sClient.BootstrapTykOperatorSecret()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	fmt.Println("Finished bootstrapping operator secret")

	fmt.Println("Started bootstrapping portal secret")
	if appArgs.EnterprisePortalSecretEnabled {
		err = k8sClient.BootstrapTykEnterprisePortalSecret()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	fmt.Println("Started bootstrapping portal with requests to dashboard")
	if appArgs.BootstrapPortal {
		err = tykSvc.BoostrapPortal()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// restarting the dashboard to apply the new cname
	if err = k8sClient.RestartDashboardDeployment(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Finished bootstrapping portal")
}
