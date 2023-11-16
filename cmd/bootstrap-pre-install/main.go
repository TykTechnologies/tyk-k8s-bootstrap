package main

import (
	"fmt"
	"os"
	"tyk/tyk/bootstrap/data"
	"tyk/tyk/bootstrap/preinstallation"
)

func main() {
	err := data.InitBootstrapConf()
	if err != nil {
		fmt.Printf("Failed to parse bootstrap environment variables, err: %v", err)
		os.Exit(1)
	}

	err = preinstallation.PreHookInstall()
	if err != nil {
		fmt.Printf("Failed to run pre-hook job, err: %v", err)
		os.Exit(1)
	}

	fmt.Println("Pre-Hook bootstrapping succeeded, the provided license is valid!")
}
