package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	flag    bool
	flatten bool
)

func init() {
	RootCmd.AddCommand(DependencyCmd)
	DependencyCmd.AddCommand(ResolveDependecyFileCmd)

	ResolveDependecyFileCmd.Flags().StringVarP(&tipe, "type", "t", "", "Type of ecosystem or file for parsing (required)")
	ResolveDependecyFileCmd.MarkFlagRequired("type")

	ResolveDependecyFileCmd.Flags().BoolVarP(&flatten, "flatten", "f", false, "Return the list in a flattened array")
	ResolveDependecyFileCmd.Flags().BoolVarP(&flag, "flag", "", false, "feature flag")
}

// DependencyCmd - Container for holding dep root and secondary commands
var DependencyCmd = &cobra.Command{
	Use:   "dependency",
	Short: "Dependency resource",
	Long:  `Dependency resource - access data relating to dependencies and their associations`,
}

// ResolveDependecyFileCmd - Resolves dependencies in a file
var ResolveDependecyFileCmd = &cobra.Command{
	Use:   "resolve-dependencies-in-file [flags] PATHTODEPFILE",
	Short: "Resolve Dependencies in a file",
	Long:  `Get the data for a dependency file`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dep := args[0]
		deps, e := ion.ResolveDependenciesInFile(dep, tipe, flatten, flag, viper.GetString(secretKey))
		if e != nil {
			fmt.Println(e.Error())
		}

		PPrint(deps)
	},
}
