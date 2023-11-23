package main

import (
	"os"
	"strings"
	"tyk/tyk/bootstrap/k8s"
	"tyk/tyk/bootstrap/pkg/config"
	"tyk/tyk/bootstrap/tyk"

	"github.com/sirupsen/logrus"
)

const (
	secretTypeOperator = "Operator"
	secretTypePortal   = "Portal"
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
	log.WithField("level", level.String()).Info("Set the log level")

	k8sClient, err := k8s.NewClient(conf, log.WithField("Client", "Kubernetes"))
	if err != nil {
		exit(log, err)
	}

	if err = k8sClient.CheckIfRequiredDeploymentsAreReady(); err != nil {
		exit(log, err)
	}

	if conf.BootstrapDashboard || conf.OperatorKubernetesSecretName != "" || conf.DevPortalKubernetesSecretName != "" {
		conf.K8s.DashboardSvcUrl, err = k8sClient.DiscoverDashboardSvc()
		if err != nil {
			exit(log, err)
		}
	}

	tykSvc := tyk.NewClient(conf, log.WithField("Client", "Tyk"))
	var orgExists bool

	if conf.BootstrapDashboard {
		orgExists, err = tykSvc.OrgExists()
		if err != nil {
			exit(log, err)
		}

		if orgExists {
			log.Warnf("Organisation named %v with cname %v exists on Tyk Dashboard",
				conf.Tyk.Org.Name,
				conf.Tyk.Org.Cname,
			)
		} else {
			log.Info("Bootstrapping Tyk Dashboard")

			if err = tykSvc.CreateOrganisation(); err != nil {
				exit(log, err)
			}

			if err = tykSvc.CreateAdmin(); err != nil {
				exit(log, err)
			}
		}
	}

	if conf.BootstrapPortal {
		log.Info("Bootstrapping Tyk Classic Portal")

		orgExists, err = tykSvc.OrgExists()
		if err != nil {
			exit(log, err)
		}

		if !orgExists {
			if err = tykSvc.BootstrapClassicPortal(); err != nil {
				exit(log, err)
			}

			if err = k8sClient.RestartDashboard(); err != nil {
				exit(log, err)
			}
		}
	}

	// Common log message for organization existence
	if orgExists && (conf.DevPortalKubernetesSecretName != "" || conf.OperatorKubernetesSecretName != "") {
		log.WithFields(logrus.Fields{
			"organisationName":  conf.Tyk.Org.Name,
			"organisationCName": conf.Tyk.Org.Cname,
		}).Info("Organisation exists on Tyk. " +
			"Please provide the Organisation ID and Dashboard Access Key for Kubernetes secrets")
	}

	if conf.DevPortalKubernetesSecretName != "" {
		err = createK8sSecret(
			log, conf, k8sClient, conf.DevPortalKubernetesSecretName, "Tyk Developer Portal",
		)
		if err != nil {
			exit(log, err)
		}
	}

	if conf.OperatorKubernetesSecretName != "" {
		err = createK8sSecret(log, conf, k8sClient, conf.OperatorKubernetesSecretName, "Tyk Operator")
		if err != nil {
			exit(log, err)
		}
	}
}

func createK8sSecret(l *logrus.Logger, c *config.Config, client *k8s.Client, secretName, secretType string) error {
	fields := logrus.Fields{"secretName": secretName}

	if c.Tyk.Org.ID == "" {
		l.WithFields(fields).
			Warn("Given Organisation ID is empty, the Kubernetes secret will contain empty TYK_ORG")
	}

	if c.Tyk.Admin.Auth == "" {
		l.WithFields(fields).
			Warn("Given User Auth key is empty, the Kubernetes secret will contain empty TYK_AUTH")
	}

	l.WithFields(fields).Infof("Creating Kubernetes Secret for %s", secretType)

	switch {
	case strings.Contains(secretType, secretTypeOperator):
		if err := client.BootstrapTykOperatorSecret(); err != nil {
			return err
		}
	case strings.Contains(secretType, secretTypePortal):
		if err := client.BootstrapTykPortalSecret(); err != nil {
			return err
		}
	}

	return nil
}

func exit(log *logrus.Logger, err error) {
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
