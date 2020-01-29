package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(AnalysisCmd)
	AnalysisCmd.AddCommand(GetAnalysisCmd)
	AnalysisCmd.AddCommand(GetLatestAnalysisCmd)

	GetAnalysisCmd.Flags().StringVarP(&teamID, "team-id", "t", "", "ID of the team for the project (required)")
	GetAnalysisCmd.MarkFlagRequired("team-id")
	GetAnalysisCmd.Flags().StringVarP(&projectID, "project-id", "p", "", "ID of the project (required)")
	GetAnalysisCmd.MarkFlagRequired("project-id")

	GetLatestAnalysisCmd.Flags().StringVarP(&teamID, "team-id", "t", "", "ID of the team for the project (required)")
	GetLatestAnalysisCmd.MarkFlagRequired("team-id")
	GetLatestAnalysisCmd.Flags().StringVarP(&projectID, "project-id", "p", "", "ID of the project (required)")
	GetLatestAnalysisCmd.MarkFlagRequired("project-id")

}

// AnalysisCmd - Container for holding project root and secondary commands
var AnalysisCmd = &cobra.Command{
	Use:   "analysis",
	Short: "Analysis resource",
	Long:  `Analysis resource - access data relating to analyses and their associations`,
}

// GetAnalysisCmd - Container for holding project root and secondary commands
var GetAnalysisCmd = &cobra.Command{
	Use:   "get-analysis [flags] ANALYSISID",
	Short: "Get Analysis",
	Long:  `Get the data for an analysis`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		a, e := ion.GetAnalysis(id, teamID, projectID, viper.GetString(secretKey))
		if e != nil {
			fmt.Println(e.Error())
		}

		PPrint(a)
	},
}

// GetLatestAnalysisCmd - Container for holding project root and secondary commands
var GetLatestAnalysisCmd = &cobra.Command{
	Use:   "get-latest-analysis [flags]",
	Short: "Get Latest Analysis",
	Long:  `Get the data for a latest analysis`,
	Run: func(cmd *cobra.Command, args []string) {
		a, e := ion.GetLatestAnalysis(teamID, projectID, viper.GetString(secretKey))
		if e != nil {
			fmt.Println(e.Error())
		}

		PPrint(a)
	},
}
