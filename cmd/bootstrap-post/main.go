package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"tyk/tyk/bootstrap/data"
	"tyk/tyk/bootstrap/k8s"
	"tyk/tyk/bootstrap/tyk"
)

func main() {
	k8sClient, err := k8s.NewClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	k8sClient.AppArgs, err = data.InitAppDataPostInstall()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if k8sClient.AppArgs.BootstrapDashboard {
		if err = k8sClient.DiscoverDashboardSvc(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		k8sClient.AppArgs.DashboardSvcAddr = fmt.Sprintf("%s://%s.%s.svc.cluster.local:%d",
			k8sClient.AppArgs.DashboardSvcProto,
			k8sClient.AppArgs.DashboardSvcName,
			k8sClient.AppArgs.ReleaseNamespace,
			k8sClient.AppArgs.DashboardSvcPort,
		)
	}

	err = k8sClient.CheckIfRequiredDeploymentsAreReady()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tp := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: k8sClient.AppArgs.InsecureSkipVerify},
	}
	client := http.Client{Transport: tp}

	tykSvc := tyk.NewTykService(client, k8sClient.AppArgs)

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
	if k8sClient.AppArgs.OperatorSecretEnabled {
		err = k8sClient.BootstrapTykOperatorSecret()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	fmt.Println("Finished bootstrapping operator secret")

	fmt.Println("Started bootstrapping portal secret")
	if k8sClient.AppArgs.EnterprisePortalSecretEnabled {
		err = k8sClient.BootstrapTykEnterprisePortalSecret()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	fmt.Println("Started bootstrapping portal with requests to dashboard")
	if k8sClient.AppArgs.BootstrapClassicPortal {
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
