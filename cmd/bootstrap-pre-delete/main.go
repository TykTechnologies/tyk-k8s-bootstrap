package main

import (
	"fmt"
	"os"
	"tyk/tyk/bootstrap/k8s"
	"tyk/tyk/bootstrap/tyk/data"
)

func main() {
	k8sClient, err := k8s.NewClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	k8sClient.AppArgs = data.InitAppDataPreDelete()

	err = k8sClient.ExecutePreDeleteOperations()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
