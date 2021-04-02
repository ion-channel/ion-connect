package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ion-channel/ionic"
	"github.com/ion-channel/ionic/pagination"
	"github.com/ion-channel/ionic/projects"
	ionicspdx "github.com/ion-channel/ionic/spdx"
	"github.com/spdx/tools-golang/tvloader"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	spdxVersion string
	rulesetID   string
	packageName string
	pocEmail    string
	pocName     string
)

func init() {
	RootCmd.AddCommand(ProjectCmd)
	ProjectCmd.AddCommand(AddAliasCmd)
	ProjectCmd.AddCommand(GetProjectCmd)
	ProjectCmd.AddCommand(GetProjectsCmd)
	ProjectCmd.AddCommand(CreateProjectsCSVCmd)
	ProjectCmd.AddCommand(CreateProjectCmd)
	ProjectCmd.AddCommand(CreateProjectsSPDXCmd)
	ProjectCmd.AddCommand(SetSourceCmd)
	ProjectCmd.AddCommand(SetTypeCmd)

	SetSourceCmd.Flags().StringVarP(&projectID, "project-id", "p", "", "ID of the project (required)")
	SetSourceCmd.MarkFlagRequired("project-id")

	SetTypeCmd.Flags().StringVarP(&projectID, "project-id", "p", "", "ID of the project (required)")
	SetTypeCmd.MarkFlagRequired("project-id")
	SetTypeCmd.Flags().StringVarP(&source, "source-location", "s", "", "Source location of project (required except for source_unavailable type)")
	SetTypeCmd.Flags().StringVarP(&branch, "source-branch-name", "b", "", "Source branch of project")

	AddAliasCmd.Flags().StringVarP(&teamID, "team-id", "t", "", "ID of the team for the project (required)")
	AddAliasCmd.MarkFlagRequired("team-id")
	AddAliasCmd.Flags().StringVarP(&projectID, "project-id", "p", "", "ID of the project (required)")
	AddAliasCmd.MarkFlagRequired("project-id")

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

	CreateProjectsSPDXCmd.Flags().StringVarP(&spdxVersion, "spdx-version", "", "2.1", "SPDX version 2.1 or 2.2 to import")
	CreateProjectsSPDXCmd.Flags().StringVarP(&teamID, "team-id", "t", "", "ID of the team for the project (required)")
	CreateProjectsSPDXCmd.MarkFlagRequired("team-id")
	CreateProjectsSPDXCmd.Flags().StringVarP(&rulesetID, "ruleset-id", "r", "", "ID of the ruleset for the project (required)")
	CreateProjectsSPDXCmd.MarkFlagRequired("ruleset-id")
	CreateProjectsSPDXCmd.Flags().StringVarP(&packageName, "package-name", "", "", "package name of the project to add (must be present in SPDX file)")
	CreateProjectsSPDXCmd.Flags().StringVarP(&pocEmail, "poc-email", "", "", "Point of Contact (PoC) email to be used for the project")
	CreateProjectsSPDXCmd.Flags().StringVarP(&pocName, "poc-name", "", "", "Point of Contact (PoC) name to be used for the project")
}

// ProjectCmd - Container for holding project root and secondary commands
var ProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Project resource",
	Long:  `Project resource - access data relating to projects and their associations`,
}

// SetSourceCmd - Set the source location of a project
var SetSourceCmd = &cobra.Command{
	Use:   "set-source LOCATION",
	Short: "Set source",
	Long:  "Set the source for the project",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		Src := args[0]
		p, e := ion.GetProject(projectID, teamID, viper.GetString(secretKey))
		if e != nil {
			fmt.Println(e.Error())
			return
		}

		p.Source = &Src
		p, e = ion.UpdateProject(p, viper.GetString(secretKey))
		if e != nil {
			fmt.Println(e.Error())
		}
		PPrint(p)
	},
}

// SetTypeCmd - set project source type (git, http, etc...)
var SetTypeCmd = &cobra.Command{
	Use:   "set-type TYPE [SOURCE] [BRANCH]",
	Short: "Set type",
	Long:  "Set the type for the project",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		Type := args[0]
		if Type != "source_unavailable" && len(source) == 0 {
			fmt.Println("Source location required")
			return
		}
		p, e := ion.GetProject(projectID, teamID, viper.GetString(secretKey))
		if e != nil {
			fmt.Println(e.Error())
			return
		}

		p.Type = &Type
		if *p.Type == "source_unavailable" {
			// source unavailable means the source needs to be blank
			empty := ""
			p.Source = &empty
		} else if len(source) != 0 {
			// other project types can set the source if present
			p.Source = &source
		}
		if len(branch) != 0 {
			p.Branch = &branch
		}

		p, e = ion.UpdateProject(p, viper.GetString(secretKey))
		if e != nil {
			fmt.Println(e.Error())
		}
		PPrint(p)
	},
}

// AddAliasCmd - Container for holder add alias cmd
var AddAliasCmd = &cobra.Command{
	Use:   "add-alias NAME VERSION [ORG]",
	Short: "Add Alias",
	Long:  `Add an alias to a project`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		options := ionic.AddAliasOptions{
			Name:      args[0],
			ProjectID: projectID,
			TeamID:    teamID,
			Version:   args[1],
		}
		if len(args) == 3 {
			options.Org = args[2]
		}
		p, e := ion.AddAlias(options, viper.GetString(secretKey))
		if e != nil {
			fmt.Println(e.Error())
		}

		PPrint(p)
	},
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
		page := pagination.New(offset, limit)
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

// CreateProjectsSPDXCmd - Attempts to create a project from an SPDX V2.1 or V2.2 file
var CreateProjectsSPDXCmd = &cobra.Command{
	Use:   "create-projects-spdx [flags] PATHTOSPDX",
	Short: "Create Projects SPDX",
	Long:  `Create projects from an spdx file`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// open the SPDX file
		spdxFile := args[0]
		r, err := os.Open(spdxFile)
		if err != nil {
			fmt.Printf("Error while opening %v for reading: %v", spdxFile, err)
			return
		}
		defer r.Close()

		// check the version of SPDX file
		if len(args) > 1 {
			spdxVersion = args[1]
		}

		// if we got here, the file is now loaded into memory.
		fmt.Printf("Read %s\n\n", spdxFile)

		// projects to create
		projs := make([]projects.Project, 0)

		// parse on SPDX 2.1, or 2.2. Default to 2.2 if not provided
		if spdxVersion == "2.1" {
			// try to load the SPDX file's contents as a tag-value file, version 2.1
			doc, err := tvloader.Load2_1(r)
			if err != nil {
				fmt.Printf("Could not load %v. Error from SPDX library: %v\n", spdxFile, err)
				printSPDXErrorHelp(err)
				return
			}

			// All project packages
			packageProjs, err := ionicspdx.ProjectPackageFromSPDX2_1(doc, packageName)
			if err != nil {
				fmt.Printf("Failed finding project packages %v: %v. Skipping project packages.", spdxFile, err)
			} else {
				projs = append(projs, packageProjs...)
			}

			// Top level SPDX Document Creation Information (DocumentInfo)
			p, err := ionicspdx.ProjectFromSPDX2_1(doc)
			if err != nil {
				fmt.Printf("Failed finding documentinfo %v: %v. Skipping creating project from documentinfo.", spdxFile, err)
			} else {
				projs = append(projs, p)
			}

		} else if spdxVersion == "2.2" || spdxVersion == "" {
			// try to load the SPDX file's contents as a tag-value file, version 2.2
			doc, err := tvloader.Load2_2(r)
			if err != nil {
				fmt.Printf("Could not load %v. Error from SPDX library: %v\n", spdxFile, err)
				printSPDXErrorHelp(err)
				return
			}
			// All project packages
			packageProjs, err := ionicspdx.ProjectPackageFromSPDX2_2(doc, packageName)
			if err != nil {
				fmt.Printf("Failed finding project packages %v: %v. Skipping project packages.", spdxFile, err)
			} else {
				projs = append(projs, packageProjs...)
			}

			// Top level SPDX Document Creation Information (DocumentInfo)
			p, err := ionicspdx.ProjectFromSPDX2_2(doc)
			if err != nil {
				fmt.Printf("Failed finding documentinfo %v: %v. Skipping creating project from documentinfo.", spdxFile, err)
			} else {
				projs = append(projs, p)
			}

		}

		if len(projs) == 0 {
			printSPDXErrorHelp(fmt.Errorf("no packages found in SPDX file"))
		}

		for i := range projs {
			projs[i].RulesetID = &rulesetID
			projs[i].TeamID = &teamID

			// add name and email if supplied
			if pocName != "" {
				projs[i].POCName = pocName
			}

			if pocEmail != "" {
				projs[i].POCEmail = pocEmail
			}

			errs, err := projs[i].Validate(cc, uu, viper.GetString(secretKey))
			if err != nil {
				fmt.Printf("\nProject %v doesn't meet Ion requirements: %v. Details: \n", *(projs[i].Name), err.Error())
				for name, e := range errs {
					fmt.Printf("%v : %v\n", name, e)
				}
				fmt.Printf("\nStopping import.")
				return
			}

			projs[i].Active = true

			fmt.Printf("\nProject we're creating: %+v\n", projs[i])
			res, err := ion.CreateProject(&projs[i], teamID, viper.GetString(secretKey))
			if err != nil {
				fmt.Printf("\nCouldn't create project: %v\n", err.Error())
				fmt.Printf("Continuing on to create other projects.\n")
				continue
			}

			PPrint(res)
		}

	},
}

func printSPDXErrorHelp(e error) {
	fmt.Printf("\n*******************************")
	fmt.Printf("\nThe error '%v' prevented the SPDX file from being parsed.\n", e.Error())
	fmt.Printf("\nPlease check the SPDX file contents follow SPDX 2.2 or 2.1 specifications:\n")
	fmt.Printf("SPDX 2.2: https://spdx.github.io/spdx-spec/\n")
	fmt.Printf("SPDX 2.1: https://spdx.dev/spdx-specification-21-web-version/\n")
	fmt.Printf("\nThis online tool can help identify issues: https://tools.spdx.org/app/validate/\n")
	fmt.Printf("*******************************\n")
}
