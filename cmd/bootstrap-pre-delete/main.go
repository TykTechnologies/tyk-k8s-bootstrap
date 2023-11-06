package main

import (
	"fmt"
	"os"
	"tyk/tyk/bootstrap/k8s"
	"tyk/tyk/bootstrap/pkg/config"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	k8sClient, err := k8s.NewClient(conf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = k8sClient.ExecutePreDeleteOperations()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
