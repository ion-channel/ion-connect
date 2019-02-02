package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/ion-channel/ionic"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	key       = "IONCHANNEL_SECRET_KEY"
	api       = "IONCHANNEL_ENDPOINT_URL"
	bucket    = "IONCHANNEL_DROP_BUCKET"
	secretKey = "secret_key"
)

var (
	output    io.Writer
	cfgFile   string
	ion       *ionic.IonClient
	teamID    string
	projectID string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ion-connect",
	Short: "Provides a micro level interface for performing supply chain analysis of a project",
	Long: `ion-connect is a CLI tool that allows for rich interaction with the Ion Channel API to
perform supply chain analysis for a project.
`,
}

func init() {
	output = os.Stdout

	cobra.OnInitialize(initDefaults, initEnvs, initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $PWD/.ionize.yaml)")
}

func initDefaults() {
	// viper.SetDefault("api", "https://api.ionchannel.io")
}

func initEnvs() {
	viper.BindEnv(secretKey, key)
	viper.BindEnv("api", api)
	viper.BindEnv("bucket", bucket)
}

func initConfig() {
	viper.SetConfigType("yaml")

	viper.SetConfigName("credentials")
	viper.AddConfigPath("/etc/ionchannel/")
	viper.AddConfigPath("$HOME/.ionchannel/")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil && (viper.GetString(secretKey) == "" || viper.GetString("api") == "") {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	ion, _ = ionic.New(viper.GetString("api"))
}

func init() {

	RootCmd.AddCommand(ProjectCmd)
}
