package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	yaml "gopkg.in/yaml.v2"
)

func init() {
	RootCmd.AddCommand(configureCmd)
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure for operation ion-connect",
	Long:  `Configure ion-connect for use with a config file in the user's home.`,
	Run: func(cmd *cobra.Command, args []string) {
		viper.ReadInConfig()
		creds := viper.GetString(secretKey)
		var truncCreds string
		if len(creds) > 4 {
			truncCreds = creds[len(creds)-4 : len(creds)]
		}

		fmt.Printf("Ion Channel Api Key [%s]: ", truncCreds)
		newCreds, e := terminal.ReadPassword(int(os.Stdin.Fd()))
		if e != nil {
			fmt.Println(e.Error())
			return
		}
		if len(newCreds) != 0 {
			saveCredentials(string(newCreds))
		}
	},
}

func saveCredentials(creds string) {
	credMap := make(map[string]string)
	credMap[secretKey] = creds
	yamlCreds, e := yaml.Marshal(&credMap)
	if e != nil {
		fmt.Println(e.Error())
		return
	}

	usr, e := user.Current()
	if e != nil {
		fmt.Println(e.Error())
		return
	}

	e = os.MkdirAll(strings.Replace(configPath, "$HOME", usr.HomeDir, -1), 0700)
	if e != nil {
		fmt.Println(e.Error())
		return
	}

	path := strings.Replace(configFilePath, "~", usr.HomeDir, -1)
	e = ioutil.WriteFile(path, yamlCreds, 0600)
	if e != nil {
		fmt.Println(e.Error())
	}
}
