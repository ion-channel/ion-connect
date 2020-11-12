package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ion-channel/ionic/requests"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(ReportCmd)
	ReportCmd.AddCommand(GetAnalysisReportCmd)
	ReportCmd.AddCommand(GetExportData)

	GetAnalysisReportCmd.Flags().StringVarP(&teamID, "team-id", "t", "", "ID of the team for the project (required)")
	GetAnalysisReportCmd.MarkFlagRequired("team-id")
	GetAnalysisReportCmd.Flags().StringVarP(&projectID, "project-id", "p", "", "ID of the project (required)")
	GetAnalysisReportCmd.MarkFlagRequired("project-id")

	GetExportData.Flags().BoolVar(&skel, "print", false, "Print an example export data json skeleton. ids is 1-n project ids.")

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

// GetExportData - Container for holding get exported data root and secondary commands
var GetExportData = &cobra.Command{
	Use:   "export-data [flags] PATHTOJSON",
	Short: "Export Data",
	Long:  `Export Data - requires team id and project id(s) from a json file`,
	Args: func(cmd *cobra.Command, args []string) error {
		if skel {
			return nil
		}
		if len(args) != 1 {
			return fmt.Errorf("accepts 1 arg(s), received 0")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if skel {
			fmt.Println("{")
			fmt.Println("   \"team_id\":\"\",")
			fmt.Println("   \"ids\": [")
			fmt.Println("      \"project_id\",")
			fmt.Println("      \"2nd project_id\",")
			fmt.Println("      \"...\",")
			fmt.Println("      \"last project_id\"")
			fmt.Println("   ]")
			fmt.Println("}")
			return
		}

		filename := args[0]

		f, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		var data requests.ByIDsAndTeamID
		err = json.Unmarshal(f, &data)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if data.TeamID == "" {
			fmt.Println("team_id is required")
			return
		}

		if len(data.IDs) == 0 {
			fmt.Println("one or more project ids are required in ids json filed")
			return
		}

		r, err := ion.GetExportedProjectsData(data.IDs, data.TeamID, viper.GetString(secretKey))
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		PPrint(r)
	},
}
