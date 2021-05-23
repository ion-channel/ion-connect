package cmd

import (
	"fmt"

	"github.com/ion-channel/ionic/dependencies"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	flag    bool
	flatten bool
)

func init() {
	RootCmd.AddCommand(DependencyCmd)
	DependencyCmd.AddCommand(ResolveDependencyFileCmd)
	DependencyCmd.AddCommand(GetVersionsCmd)

	ResolveDependencyFileCmd.Flags().StringVarP(&tipe, "type", "t", "", "Type of ecosystem or file for parsing (required)")
	ResolveDependencyFileCmd.MarkFlagRequired("type")

	ResolveDependencyFileCmd.Flags().BoolVarP(&flatten, "flatten", "f", false, "Return the list in a flattened array")
	ResolveDependencyFileCmd.Flags().BoolVarP(&flag, "flag", "", false, "feature flag")

	GetVersionsCmd.Flags().StringVarP(&tipe, "type", "t", "", "Type of ecosystem or file for parsing (required)")
	GetVersionsCmd.MarkFlagRequired("type")
}

// DependencyCmd - Container for holding dep root and secondary commands
var DependencyCmd = &cobra.Command{
	Use:   "dependency",
	Short: "Dependency resource",
	Long:  `Dependency resource - access data relating to dependencies and their associations`,
}

// ResolveDependencyFileCmd - Resolves dependencies in a file
var ResolveDependencyFileCmd = &cobra.Command{
	Use:   "resolve-dependencies-in-file [flags] PATHTODEPFILE",
	Short: "Resolve Dependencies in a file",
	Long:  `Get the data for a dependency file`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dep := args[0]
		o := dependencies.DependencyResolutionRequest{
			File:      dep,
			Flatten:   flatten,
			Ecosystem: tipe,
		}

		deps, e := ion.ResolveDependenciesInFile(o, viper.GetString(secretKey))
		if e != nil {
			fmt.Println(e.Error())
		}

		PPrint(deps)
	},
}

// GetVersionsCmd - Resolves versions for a dependency
var GetVersionsCmd = &cobra.Command{
	Use:   "get-versions [flags] NAME",
	Short: "Resolves versions for a dependency",
	Long:  `Get a list of versions for a dependency`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dep := args[0]

		vs, e := ion.GetVersionsForDependency(dep, tipe, viper.GetString(secretKey))
		if e != nil {
			fmt.Println(e.Error())
		}

		PPrint(vs)
	},
}
