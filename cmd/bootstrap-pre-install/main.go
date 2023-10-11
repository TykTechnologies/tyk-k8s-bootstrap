package main

import (
	"fmt"
	"os"
	"tyk/tyk/bootstrap/k8s"
)

func main() {
	err := k8s.PreHookInstall()
	if err != nil {
		fmt.Printf("Failed to run pre-hook job, err: %v", err)
		os.Exit(1)
	}

	fmt.Println("Pre-Hook bootstrapping succeeded, the provided license is valid!")
}
