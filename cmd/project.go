package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	ProjectCmd.AddCommand(GetProjectCmd)

	GetProjectCmd.Flags().StringVarP(&teamID, "team-id", "t", "", "ID of the team for the project (required)")
	GetProjectCmd.MarkFlagRequired("team-id")
	GetProjectCmd.Flags().StringVarP(&projectID, "project-id", "p", "", "ID of the project (required)")
	GetProjectCmd.MarkFlagRequired("project-id")
}

// ProjectCmd - Container for holding project root and secondary commands
var ProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Project resource",
	Long:  `Project resource - access data relating to projects and their associations`,
}

// GetProjectCmd - Container for holding project root and secondary commands
var GetProjectCmd = &cobra.Command{
	Use:   "get-project",
	Short: "Get Project",
	Long:  `Get the data for a project`,
	Run: func(cmd *cobra.Command, args []string) {
		p, e := ion.GetProject(projectID, teamID, viper.GetString(secretKey))
		if e != nil {
			fmt.Println(e.Error())
		}

		j, e := json.Marshal(p)
		if e != nil {
			fmt.Println(e.Error())
		}
		fmt.Println(string(j))
	},
}
