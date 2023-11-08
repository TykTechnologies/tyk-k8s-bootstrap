package main

import (
	"github.com/sirupsen/logrus"
	"os"
	"tyk/tyk/bootstrap/k8s"
	"tyk/tyk/bootstrap/pkg/config"
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

	if err = k8sClient.ExecutePreDeleteOperations(); err != nil {
		exit(log, err)
	}
}

func exit(log *logrus.Logger, err error) {
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
