package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ion-channel/ionic/pagination"
	"github.com/ion-channel/ionic/projects"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(ProjectCmd)
	ProjectCmd.AddCommand(GetProjectCmd)
	ProjectCmd.AddCommand(GetProjectsCmd)
	ProjectCmd.AddCommand(CreateProjectsCSVCmd)
	ProjectCmd.AddCommand(CreateProjectCmd)

	GetProjectCmd.Flags().StringVarP(&teamID, "team-id", "t", "", "ID of the team for the project (required)")
	GetProjectCmd.MarkFlagRequired("team-id")
	GetProjectCmd.Flags().StringVarP(&projectID, "project-id", "p", "", "ID of the project (required)")
	GetProjectCmd.MarkFlagRequired("project-id")

	GetProjectsCmd.Flags().StringVarP(&teamID, "team-id", "t", "", "ID of the team for the project (required)")
	GetProjectsCmd.Flags().IntVarP(&limit, "limit", "", 10, "maximum count of projects")
	GetProjectsCmd.Flags().IntVarP(&offset, "offset", "", 0, "beginning index for project set")
	GetProjectsCmd.MarkFlagRequired("team-id")

	CreateProjectsCSVCmd.Flags().StringVarP(&teamID, "team-id", "t", "", "ID of the team for the project (required)")
	CreateProjectsCSVCmd.MarkFlagRequired("team-id")

	CreateProjectCmd.Flags().BoolVar(&skel, "print", false, "Print an example create project json skeleton")
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

		PPrint(p)
	},
}

// GetProjectsCmd - Container for holding project root and secondary commands
var GetProjectsCmd = &cobra.Command{
	Use:   "get-projects",
	Short: "Get Projects",
	Long:  `Get the data for a projects in a team`,
	Run: func(cmd *cobra.Command, args []string) {
		page := pagination.New(limit, offset)
		ps, e := ion.GetProjects(teamID, viper.GetString(secretKey), page, nil)
		if e != nil {
			fmt.Println(e.Error())
		}

		PPrint(ps)
	},
}

// CreateProjectsCSVCmd - Container for holding project root and secondary commands
var CreateProjectsCSVCmd = &cobra.Command{
	Use:   "create-projects-csv [flags] PATHTOCSV",
	Short: "Create Projects",
	Long:  `Create projects from a Ion Channel CSV input file`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		csv := args[0]
		ps, e := ion.CreateProjectsFromCSV(csv, teamID, viper.GetString(secretKey))
		if e != nil {
			fmt.Println(e.Error())
		}

		PPrint(ps)
	},
}

// CreateProjectCmd - Container for holding create project root and secondary commands
var CreateProjectCmd = &cobra.Command{
	Use:   "create-project [flags] PATHTOJSON",
	Short: "Create Project",
	Long:  `Create project from a Ion Channel JSON input file`,
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
			fmt.Println("   \"ruleset_id\":\"\",")
			fmt.Println("   \"name\":\"\",")
			fmt.Println("   \"type\":\"\",")
			fmt.Println("   \"source\":\"\",")
			fmt.Println("   \"description\":\"\",")
			fmt.Println("   \"username\":\"\",")
			fmt.Println("   \"password\":\"\"")
			fmt.Println("}")
			return
		}

		filename := args[0]

		f, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		var data projects.Project
		err = json.Unmarshal(f, &data)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		errs, err := data.Validate(cc, uu, viper.GetString(secretKey))
		if err != nil {
			for name, e := range errs {
				fmt.Printf("%v : %v\n", name, e)
			}

			fmt.Println(err.Error())
			return
		}

		data.Active = true

		r, err := ion.CreateProject(&data, teamID, viper.GetString(secretKey))
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		PPrint(r)
	},
}
