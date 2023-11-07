package main

import (
	"errors"
	"fmt"
	"os"
	"tyk/tyk/bootstrap/k8s"
	"tyk/tyk/bootstrap/pkg/config"
	"tyk/tyk/bootstrap/tyk"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		exit(err)
	}

	k8sClient, err := k8s.NewClient(conf)
	if err != nil {
		exit(err)
	}

	if err = k8sClient.CheckIfRequiredDeploymentsAreReady(); err != nil {
		exit(err)
	}

	orgExists := false

	tykSvc := tyk.NewService(conf)
	if err = tykSvc.OrgExists(); err != nil {
		if !errors.Is(err, tyk.ErrOrgExists) {
			exit(err)
		}

		orgExists = true
	}

	if !orgExists {
		if err = tykSvc.CreateOrganisation(); err != nil {
			exit(err)
		}

		if err = tykSvc.CreateAdmin(); err != nil {
			exit(err)
		}

		if conf.BootstrapPortal {
			if err = tykSvc.BootstrapClassicPortal(); err != nil {
				exit(err)
			}
		}

		if err = k8sClient.RestartDashboard(); err != nil {
			exit(err)
		}
	}

	if conf.DevPortalKubernetesSecretName != "" {
		err = k8sClient.BootstrapTykPortalSecret()
		if err != nil {
			exit(err)
		}
	}

	if conf.OperatorKubernetesSecretName != "" {
		err = k8sClient.BootstrapTykOperatorSecret()
		if err != nil {
			exit(err)
		}
	}
}

func exit(err error) {
	if err != nil {
		fmt.Printf("[ERROR]: %v", err)
		os.Exit(1)
	}
}
