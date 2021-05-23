package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Resource enum for defining resource types
type Resource string

const (
	// Repo repository resource type
	Repo Resource = "repos"
	// Product product (cpe, software) resource type
	Product Resource = "products"
	// Package package (dependencies, libraries) resource type
	Package Resource = "packages"
)

var (
	resource string
)

// Valid determines if a potential resource is actually one
func (r *Resource) Valid() bool {
	switch *r {
	case Repo, Product, Package, "":
		return true
	}
	return false
}

func init() {
	RootCmd.AddCommand(SearchCmd)
	SearchCmd.Flags().StringVarP(&query, "query", "q", "", "String text to search Ion Channel for matches (required)")
	SearchCmd.Flags().StringVarP(&resource, "tbs", "t", "products", "Requests that the search is performed against the provided resource type")
	SearchCmd.MarkFlagRequired("query")
}

// SearchCmd For a query string provide a result set
// comprised of known data about a software artifact
var SearchCmd = &cobra.Command{
	Use:   "search",
	Short: "For a query string provide a result set",
	Long: `For a query string provide a result set
	comprised of known data about a software artifact`,
	Args: func(cmd *cobra.Command, args []string) error {
		if (*Resource)(&resource).Valid() {
			return nil
		}
		return fmt.Errorf("invalid resource type specified: %s", resource)
	},
	Run: func(cmd *cobra.Command, args []string) {
		r, _, e := ion.GetSearch(query, resource, viper.GetString(secretKey))
		if e != nil {
			fmt.Println(e.Error())
		}

		PPrint(r)
	},
}
