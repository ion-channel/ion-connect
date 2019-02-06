package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(ReportCmd)
	ReportCmd.AddCommand(GetAnalysisReportCmd)

	GetAnalysisReportCmd.Flags().StringVarP(&teamID, "team-id", "t", "", "ID of the team for the project (required)")
	GetAnalysisReportCmd.MarkFlagRequired("team-id")
	GetAnalysisReportCmd.Flags().StringVarP(&projectID, "project-id", "p", "", "ID of the project (required)")
	GetAnalysisReportCmd.MarkFlagRequired("project-id")

}

// ReportCmd - Container for holding project root and secondary commands
var ReportCmd = &cobra.Command{
	Use:   "report",
	Short: "Report resource",
	Long:  `Report resource - access data relating to analyses, evaluations and their associations`,
}

// GetAnalysisReportCmd - Container for holding project root and secondary commands
var GetAnalysisReportCmd = &cobra.Command{
	Use:   "get-analysis [flags] ANALYSISID",
	Short: "Get Analysis",
	Long:  `Get the data for an analysis`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		a, e := ion.GetAnalysisReport(id, teamID, projectID, viper.GetString(secretKey))
		if e != nil {
			fmt.Println(e.Error())
		}

		PPrint(a)
	},
}
