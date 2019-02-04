package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(ScannerCmd)
	ScannerCmd.AddCommand(GetAnalysisStatusCmd)

	GetAnalysisStatusCmd.Flags().StringVarP(&teamID, "team-id", "t", "", "ID of the team for the project (required)")
	GetAnalysisStatusCmd.MarkFlagRequired("team-id")
	GetAnalysisStatusCmd.Flags().StringVarP(&projectID, "project-id", "p", "", "ID of the project (required)")
	GetAnalysisStatusCmd.MarkFlagRequired("project-id")

}

// ScannerCmd - Container for holding project root and secondary commands
var ScannerCmd = &cobra.Command{
	Use:   "scanner",
	Short: "Scanner resource",
	Long:  `Scanner resource - access data relating to the status of analyses and their associations`,
}

// GetAnalysisStatusCmd - Container for holding project root and secondary commands
var GetAnalysisStatusCmd = &cobra.Command{
	Use:   "get-analysis-status [flags] ANALYSISID",
	Short: "Get Scanner",
	Long:  `Get the data for a project`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		a, e := ion.GetAnalysisStatus(id, teamID, projectID, viper.GetString(secretKey))
		if e != nil {
			fmt.Println(e.Error())
		}

		PPrint(a)
	},
}
