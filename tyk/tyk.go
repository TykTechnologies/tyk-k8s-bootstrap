package tyk

import (
	"crypto/tls"
	"github.com/sirupsen/logrus"
	"net/http"
	"tyk/tyk/bootstrap/pkg/config"
)

type Service struct {
	httpClient http.Client
	appArgs    *config.Config
	l          *logrus.Entry
}

// NewService returns a new service to interact with Tyk.
func NewService(args *config.Config, l *logrus.Entry) Service {
	tp := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: args.InsecureSkipVerify},
	}

	return Service{httpClient: http.Client{Transport: tp}, appArgs: args, l: l}
}
