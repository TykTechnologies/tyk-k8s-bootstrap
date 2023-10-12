package tyk

import (
	"net/http"
	"tyk/tyk/bootstrap/data"
)

type Service struct {
	httpClient http.Client
	appArgs    *data.AppArguments
}

func NewTykService(c http.Client, args *data.AppArguments) Service {
	return Service{httpClient: c, appArgs: args}
}
