package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var (
	// Version of ion-connect.
	Version string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of ion-connect",
	Long:  `All software has versions. This is mine`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ion-connect %v\n", Version)
	},
}
