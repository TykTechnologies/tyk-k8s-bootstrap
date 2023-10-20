package main

import (
	"fmt"
	"github.com/TykTechnologies/tyk-k8s-bootstrap/data"
	"github.com/TykTechnologies/tyk-k8s-bootstrap/predelete"
	"os"
)

func main() {
	err := data.InitAppDataPreDelete()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = predelete.ExecutePreDeleteOperations()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
