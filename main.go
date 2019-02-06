package main

import (
	"fmt"
	"os"

	"github.com/ion-channel/ion-connect/cmd"
)

var (
	appVersion string
)

func main() {
	cmd.Version = appVersion
	if appVersion == "" {
		cmd.Version = "local-build"
	}

	err := cmd.RootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
