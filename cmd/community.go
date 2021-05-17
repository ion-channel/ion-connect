package cmd

import (
	"errors"
	"fmt"

	"github.com/ion-channel/ionic"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	comparands []string
	byactor    bool
)

func init() {
	RootCmd.AddCommand(CommunityCmd)
	CommunityCmd.AddCommand(GetRepoCmd)
	CommunityCmd.AddCommand(GetReposForActorCmd)
	CommunityCmd.AddCommand(GetReposInCommonCmd)
	GetReposInCommonCmd.Flags().StringSliceVarP(&comparands, "comparands", "c",
		[]string{}, "Set of repository names for comparison in common. Example: -c org/repo")
	GetReposInCommonCmd.MarkFlagRequired("comparands")
	GetReposInCommonCmd.Flags().BoolVarP(&byactor, "by-actor", "b",
		false, "Request relations using actors instead of committers")
}

// CommunityCmd - Container for holding com root and secondary commands
var CommunityCmd = &cobra.Command{
	Use:   "community",
	Short: "Community resource",
	Long:  `Community resource - access data relating to communities, repositories and their associations`,
}

// GetRepoCmd - Container for holding get repo command
var GetRepoCmd = &cobra.Command{
	Use:   "get-repo NAME",
	Short: "Get a repository",
	Long:  `Get the data for repository in a community`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("name is a required argument")
		}
		name = args[0]
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		r, e := ion.GetRepo(name, viper.GetString(secretKey))
		if e != nil {
			fmt.Println(e.Error())
		}

		PPrint(r)
	},
}

// GetReposForActorCmd - Container for holding get repos for actor command
var GetReposForActorCmd = &cobra.Command{
	Use:   "get-repos-for-actor NAME",
	Short: "Get repos for a community actor",
	Long:  `Get the data for repositories connected to an actor in a community`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("name is a required argument")
		}
		name = args[0]
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ps, e := ion.GetReposForActor(name, viper.GetString(secretKey))
		if e != nil {
			fmt.Println(e.Error())
		}

		PPrint(ps)
	},
}

// GetReposInCommonCmd - Container for holding get repos in common command
var GetReposInCommonCmd = &cobra.Command{
	Use:   "get-repos-in-common SUBJECT",
	Short: "Get repos in common data",
	Long: `Get the data for repositories connected to another
	 repository through actors in a community`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("subject is a required argument")
		}
		name = args[0]
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		options := ionic.GetReposInCommonOptions{
			Subject:    name,
			Comparands: comparands,
			ByActor:    byactor,
		}
		ps, e := ion.GetReposInCommon(options, viper.GetString(secretKey))
		if e != nil {
			fmt.Println(e.Error())
		}

		PPrint(ps)
	},
}
