package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ion-channel/ionic"
	"github.com/ion-channel/ionic/pagination"
	"github.com/ion-channel/ionic/projects"
	"github.com/ion-channel/ionic/spdx"
	"github.com/spdx/tools-golang/tvloader"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	spdxVersion         string
	rulesetID           string
	pocEmail            string
	pocName             string
	includeDependencies bool
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

	CreateProjectsSPDXCmd.Flags().StringVarP(&spdxVersion, "spdx-version", "", "2.1", "SPDX version 2.1 or 2.2 to import (default: 2.1)")
	CreateProjectsSPDXCmd.Flags().StringVarP(&teamID, "team-id", "t", "", "ID of the team for the project (required)")
	CreateProjectsSPDXCmd.MarkFlagRequired("team-id")
	CreateProjectsSPDXCmd.Flags().StringVarP(&rulesetID, "ruleset-id", "r", "", "ID of the ruleset for the project (required)")
	CreateProjectsSPDXCmd.MarkFlagRequired("ruleset-id")
	CreateProjectsSPDXCmd.Flags().StringVarP(&pocEmail, "poc-email", "", "", "Point of Contact (PoC) email to be used for the project")
	CreateProjectsSPDXCmd.Flags().StringVarP(&pocName, "poc-name", "", "", "Point of Contact (PoC) name to be used for the project")
	CreateProjectsSPDXCmd.Flags().BoolVarP(&includeDependencies, "include-dependencies", "d", true, "True if dependency packages should be imported as projects (default: false)")
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
		var err error

		spdxFilePath := args[0]
		spdxFile, err := os.Open(spdxFilePath)
		if err != nil {
			fmt.Printf("Error while opening %v for reading: %v\n", spdxFilePath, err)
			return
		}
		defer spdxFile.Close()

		// check the version of SPDX file
		if len(args) > 1 {
			spdxVersion = args[1]
		}

		var spdxDoc interface{} // store doc as version-agnostic interface
		switch spdxVersion {
		case "2.1":
			spdxDoc, err = tvloader.Load2_1(spdxFile)
			if err != nil {
				fmt.Printf("Could not load %s. Error from SPDX library: %s\n", spdxFilePath, err.Error())
				printSPDXErrorHelp(err)
				return
			}
		case "2.2":
			// try to load the SPDX file's contents as a tag-value file, version 2.2
			spdxDoc, err = tvloader.Load2_2(spdxFile)
			if err != nil {
				fmt.Printf("Could not load %s. Error from SPDX library: %s\n", spdxFilePath, err.Error())
				printSPDXErrorHelp(err)
				return
			}
		default:
			fmt.Printf("Invalid SPDX specification version: %s\n", spdxVersion)
			return
		}

		newProjects, err := spdx.ProjectsFromSPDX(spdxDoc, includeDependencies)
		if err != nil {
			fmt.Printf("Failed to convert SPDX packages to projects: %s\n", err.Error())
		}

		numProjects := len(newProjects)
		if numProjects == 0 {
			printSPDXErrorHelp(fmt.Errorf("no packages found in SPDX file"))
			return
		}

		fmt.Printf("Found %d projects in SBOM.\n", numProjects)

		// compare new projects to existing projects to exclude duplicates
		var projectsToCreate []projects.Project
		existingProjects, err := ion.GetProjects(teamID, viper.GetString(secretKey), pagination.AllItems, nil)
		if err != nil {
			fmt.Printf("Failed to retrieve team's existing projects to check for duplicates: %s\n", err.Error())
			projectsToCreate = newProjects
		} else {
			existingProjectNames := make(map[string]bool)
			for _, existingProject := range existingProjects {
				if existingProject.Name != nil {
					existingProjectNames[*existingProject.Name] = true
				}
			}

			numDuplicates := 0
			for _, newProject := range newProjects {
				projectAlreadyExists := existingProjectNames[*newProject.Name]
				if !projectAlreadyExists {
					projectsToCreate = append(projectsToCreate, newProject)
				} else {
					numDuplicates += 1
				}
			}

			if numDuplicates > 0 {
				numProjects = numProjects - numDuplicates
				fmt.Printf("Ignoring %d projects that already exist.\n", numDuplicates)
			}
		}

		fmt.Printf("Importing %d projects...\n", numProjects)

		// keep track of errors as we import projects
		successes := 0
		projectsErrored := make(map[string]string)
		for ii := range projectsToCreate {
			projectsToCreate[ii].RulesetID = &rulesetID
			projectsToCreate[ii].TeamID = &teamID

			if pocName != "" {
				projectsToCreate[ii].POCName = pocName
			}

			if pocEmail != "" {
				projectsToCreate[ii].POCEmail = pocEmail
			}

			errs, err := projectsToCreate[ii].Validate(cc, uu, viper.GetString(secretKey))
			if err != nil {
				fmt.Printf("[%d/%d]\tProject %s doesn't meet Ion requirements: %s. Details: \n", ii + 1, numProjects, *(projectsToCreate[ii].Name), err.Error())
				errorStored := ""
				for name, e := range errs {
					fmt.Printf("%s : %s\n", name, e)
					errorStored += fmt.Sprintf("%s : %s\n", name, e)
				}
				projectsErrored[*(projectsToCreate[ii].Name)] = errorStored
				continue
			}

			fmt.Printf("[%d/%d]\tCreating project: %s\n", ii + 1, numProjects, *projectsToCreate[ii].Name)
			_, err = ion.CreateProject(&projectsToCreate[ii], teamID, viper.GetString(secretKey))
			if err != nil {
				fmt.Printf("\tError: %s\n", err.Error())
				projectsErrored[*(projectsToCreate[ii].Name)] = err.Error()
				continue
			}

			successes += 1
		}

		fmt.Printf("Successfully created %d/%d projects.\n", successes, numProjects)

		numErrors := len(projectsErrored)
		if numErrors > 0 {
			fmt.Printf("%d errors:\n", numErrors)
			for k, v := range projectsErrored {
				fmt.Printf("%v: %v\n", k, v)
			}
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
