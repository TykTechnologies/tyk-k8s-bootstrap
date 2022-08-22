package main

import (
	"fmt"
	"os"
	"tyk/tyk/bootstrap/data"
	"tyk/tyk/bootstrap/predelete"
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
