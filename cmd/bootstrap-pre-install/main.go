package main

import (
	"fmt"
	"os"
	"tyk/tyk/bootstrap/pkg"
	"tyk/tyk/bootstrap/pkg/config"

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

	licenseIsValid, err := pkg.ValidateDashboardLicense(conf.Tyk.DashboardLicense)
	if err != nil {
		exit(log, err)
	}

	if !licenseIsValid {
		exit(log, fmt.Errorf("provided license is invalid"))
	}

	log.Info("Pre-install Hook bootstrapping succeeded, the provided license is valid!")
}

func exit(log *logrus.Logger, err error) {
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
