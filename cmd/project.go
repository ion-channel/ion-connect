package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	fmt.Println("Adding command")
	ProjectCmd.AddCommand(GetProjectCmd)
	RootCmd.AddCommand(ProjectCmd)
}

// ProjectCmd - Container for holding project root and secondary commands
var ProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Project anything to the screen",
	Long: `project is for printing anything back to the screen.
For many years people have printed back to the screen.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Inside ProjectCmd Run with args: %v\n", args)
	},
}

// GetProjectCmd - Container for holding project root and secondary commands
var GetProjectCmd = &cobra.Command{
	Use:   "get-project",
	Short: "Project anything to the screen",
	Long: `project is for printing anything back to the screen.
For many years people have printed back to the screen.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Inside GetProjectCmd Run with args: %v\n", args)
	},
}
