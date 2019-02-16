package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ion-channel/ionic"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	key            = "IONCHANNEL_SECRET_KEY"
	api            = "IONCHANNEL_ENDPOINT_URL"
	bucket         = "IONCHANNEL_DROP_BUCKET"
	secretKey      = "secret_key"
	configPath     = "$HOME/.ionchannel/"
	configFilePath = "~/.ionchannel/credentials.yaml"
)

var (
	output    io.Writer
	cfgFile   string
	ion       *ionic.IonClient
	teamID    string
	projectID string
	branch    string
	query     string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ion-connect",
	Short: "Provides a micro level interface for performing supply chain analysis of a project",
	Long: `ion-connect is a CLI tool that allows for rich interaction with the Ion Channel API to
perform supply chain analysis for a project.
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		viper.SetConfigType("yaml")
		viper.SetConfigName("credentials")
		viper.AddConfigPath(configPath)
		// The configure command will be run before there is config
		if !strings.Contains(cmd.CommandPath(), "configure") {
			initConfig()
		}
	},
}

func init() {
	output = os.Stdout

	cobra.OnInitialize(initDefaults, initEnvs)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is %s)", configFilePath))
}

func initDefaults() {
	viper.SetDefault("api", "https://api.ionchannel.io/")
}

func initEnvs() {
	viper.BindEnv(secretKey, key)
	viper.BindEnv("api", api)
	viper.BindEnv("bucket", bucket)
}

func initConfig() {
	err := viper.ReadInConfig()
	if err != nil && (viper.GetString(secretKey) == "" || viper.GetString("api") == "") {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	ion, _ = ionic.New(viper.GetString("api"))
}
