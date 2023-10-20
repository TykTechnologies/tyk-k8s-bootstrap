package tyk

import (
	"net/http"
	"tyk/tyk/bootstrap/data"
)

type Service struct {
	httpClient http.Client
	appArgs    *data.BootstrapConf
}

func NewTykService(c http.Client, args *data.BootstrapConf) Service {
	return Service{httpClient: c, appArgs: args}
}
