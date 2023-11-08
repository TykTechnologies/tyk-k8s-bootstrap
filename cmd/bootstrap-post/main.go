package main

import (
	"errors"
	"os"
	"strings"
	"tyk/tyk/bootstrap/k8s"
	"tyk/tyk/bootstrap/pkg/config"
	"tyk/tyk/bootstrap/tyk"

	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()

	conf, err := config.NewConfig()
	if err != nil {
		exit(log, err)
	}

	level, err := logrus.ParseLevel(conf.Log)
	if err != nil {
		log.Infof(
			"Failed to parse log level %v, continuing with the default log level Info, err %v", conf.Log, err,
		)

		level = logrus.InfoLevel
	}

	log.SetLevel(level)
	log.WithField("level", level.String()).Debug("Set the log level")

	k8sClient, err := k8s.NewClient(conf, log.WithField("Client", "Kubernetes"))
	if err != nil {
		exit(log, err)
	}

	if err = k8sClient.CheckIfRequiredDeploymentsAreReady(); err != nil {
		exit(log, err)
	}

	orgExists := false

	tykSvc := tyk.NewService(conf, log.WithField("Client", "Tyk"))
	if err = tykSvc.OrgExists(); err != nil {
		if !errors.Is(err, tyk.ErrOrgExists) {
			exit(log, err)
		}

		orgExists = true
	}

	if !orgExists {
		if conf.BootstrapDashboard {
			log.Info("Bootstrapping Tyk Dashboard")

			if err = tykSvc.CreateOrganisation(); err != nil {
				exit(log, err)
			}

			if err = tykSvc.CreateAdmin(); err != nil {
				exit(log, err)
			}
		}

		if conf.BootstrapPortal {
			log.Info("Bootstrapping Tyk Classic Portal")

			if err = tykSvc.BootstrapClassicPortal(); err != nil {
				exit(log, err)
			}

			if err = k8sClient.RestartDashboard(); err != nil {
				exit(log, err)
			}
		}
	}

	// Common log message for organization existence
	if orgExists && conf.DevPortalKubernetesSecretName != "" || conf.OperatorKubernetesSecretName != "" {
		log.WithFields(logrus.Fields{
			"organisationName":  conf.Tyk.Org.Name,
			"organisationCName": conf.Tyk.Org.Cname,
		}).Info("Organisation exists on Tyk. " +
			"Please provide the Organisation ID and Dashboard Access Key for Kubernetes secrets")
	}

	createK8sSecret := func(k8sClient *k8s.Client, secretName, secretType string) {
		fields := logrus.Fields{"secretName": secretName}

		if conf.Tyk.Org.ID == "" {
			log.WithFields(fields).
				Warn("Given Organisation ID is empty, the Kubernetes secret will contain empty TYK_ORG")
		}

		if conf.Tyk.Admin.Auth == "" {
			log.WithFields(fields).
				Warn("Given User Auth key is empty, the Kubernetes secret will contain empty TYK_AUTH")
		}

		log.WithFields(fields).Infof("Creating Kubernetes Secret for %s", secretType)

		switch {
		case strings.Contains(secretType, "Operator"):
			if err := k8sClient.BootstrapTykOperatorSecret(); err != nil {
				exit(log, err)
			}
		case strings.Contains(secretType, "Portal"):
			if err := k8sClient.BootstrapTykPortalSecret(); err != nil {
				exit(log, err)
			}
		}
	}

	if conf.DevPortalKubernetesSecretName != "" {
		createK8sSecret(k8sClient, conf.DevPortalKubernetesSecretName, "Tyk Developer Portal")
	}

	if conf.OperatorKubernetesSecretName != "" {
		createK8sSecret(k8sClient, conf.OperatorKubernetesSecretName, "Tyk Operator")
	}
}

func exit(log *logrus.Logger, err error) {
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
