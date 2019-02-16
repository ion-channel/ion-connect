package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(SearchCmd)
	SearchCmd.Flags().StringVarP(&query, "query", "q", "", "String text to search Ion Channel for matches (required)")
	SearchCmd.MarkFlagRequired("query")
}

// SearchCmd For a queary string provide a resutl set
// comprised of known data about a software artifact
var SearchCmd = &cobra.Command{
	Use:   "search",
	Short: "For a query string provide a result set",
	Long: `For a queary string provide a resutl set
	comprised of known data about a software artifact`,
	Run: func(cmd *cobra.Command, args []string) {
		r, e := ion.GetSearch(query, viper.GetString(secretKey))
		if e != nil {
			fmt.Println(e.Error())
		}

		PPrint(r)
	},
}
