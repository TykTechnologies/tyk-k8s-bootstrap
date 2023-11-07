package tyk

import (
	"crypto/tls"
	"net/http"
	"tyk/tyk/bootstrap/pkg/config"
)

type Service struct {
	httpClient http.Client
	appArgs    *config.Config
}

// NewService returns a new service to interact with Tyk.
func NewService(args *config.Config) Service {
	tp := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: args.InsecureSkipVerify},
	}

	return Service{httpClient: http.Client{Transport: tp}, appArgs: args}
}
