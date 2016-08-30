// config.go
//
// Copyright [2016] [Selection Pressure]
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ionconnect

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/GeertJohan/go.rice"
	"github.com/codegangsta/cli"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/yaml.v2"
)

type Command struct {
	Name        string
	Usage       string
	Method      string
	Url         string
	Flags       []Flag
	Args        Args
	Subcommands []Command
}

func (command Command) GetArgsUsage() string {
	var buffer bytes.Buffer
	for _, arg := range command.Args {
		if len(arg.Usage) > 0 {
			if !arg.Required {
				buffer.WriteString("[")
			}
			buffer.WriteString(arg.Usage)
			if !arg.Required {
				buffer.WriteString("]")
			}
			buffer.WriteString(" ")

		}
	}

	return buffer.String()
}

func (command Command) GetArgsForFlags(flagName string) Args {
	for _, flag := range command.Flags {
		if flag.Name == flagName {
			return flag.Args
		}
	}

	return []Arg{}
}

func (command Command) GetFlagsWithArgs() []Flag {
	flags := []Flag{}
	for _, flag := range command.Flags {
		if len(flag.Args) > 0 {
			flags = append(flags, flag)
		}
	}

	return flags
}

func (command Command) GetArgsUsageWithFlags(flagName string) string {
	var buffer bytes.Buffer

	for _, flag := range command.Flags {
		if flag.Name == flagName {
			for _, arg := range flag.Args {
				if len(arg.Usage) > 0 {
					if !arg.Required {
						buffer.WriteString("[")
					}
					buffer.WriteString(arg.Usage)
					if !arg.Required {
						buffer.WriteString("]")
					}
					buffer.WriteString(" ")
				}
			}
		}
	}

	return buffer.String()
}

func (args Args) GetRequiredArgsCount() int {
	var count int
	for _, arg := range args {
		Debugf("Arg (%s) is required (%b)", arg.Name, arg.Required)
		if len(arg.Usage) > 0 && arg.Required {
			count++
		}
	}

	return count
}

func (command Command) GetDefaultRequiredArgsCount() int {
	return command.Args.GetRequiredArgsCount()
}

type Arg struct {
	Name     string
	Value    string
	Usage    string
	Required bool
	Type     string
}

type Args []Arg

type Flag struct {
	Name     string
	Value    string
	Usage    string
	Type     string
	Required bool
	Args     Args
}

type Config struct {
	Version  string
	Endpoint string
	Token    string
	Commands []Command
}

func GetConfig() Config {
	configBox, err := rice.FindBox("../config")
	if err != nil {
		log.Fatal(err)
	}

	// get file contents as string
	configString, err := configBox.String("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	config := Config{}

	err = yaml.Unmarshal([]byte(configString), &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if !Test {
		config.Commands = config.Commands[:len(config.Commands)-1]
	}
	return config
}

func (config Config) FindCommandConfig(commandName string) (Command, error) {
	for _, command := range config.Commands {
		if command.Name == commandName {
			return command, nil
		}
	}

	return Command{}, errors.New("Command not found")
}

func (config Config) ProcessUrlFromConfig(commandName string, subcommandName string, params interface{}) (string, error) {
	subCommandConfig, err := config.FindSubCommandConfig(commandName, subcommandName)
	if err != nil {
		return "", err
	}

	url := subCommandConfig.Url

	templ := template.Must(template.New("url").Parse(url))
	buf := bytes.Buffer{}
	err = templ.Execute(&buf, params)
	if err != nil {
		return "", err
	}

	return string(buf.Bytes()), nil
}

func (config Config) FindSubCommandConfig(commandName string, subcommandName string) (Command, error) {
	command, err := config.FindCommandConfig(commandName)
	if err != nil {
		return Command{}, err
	}

	for _, subcommand := range command.Subcommands {
		if subcommand.Name == subcommandName {
			return subcommand, nil
		}
	}

	return Command{}, errors.New("Subcommand not found")
}

func (config Config) LoadEndpoint() string {
	endpoint := os.Getenv(ENDPOINT_ENVIRONMENT_VARIABLE)
	if endpoint == "" {
		Debugf("Endpoint env var not found returning from config file (%s)", config.Endpoint)

		return config.Endpoint
	} else {
		Debugf("Credential env var found (%s)", endpoint)
		return endpoint
	}
}

func HandleConfigure(context *cli.Context) {
	currentSecretKey := LoadCredential()
	truncatedSecretKey := currentSecretKey
	if len(currentSecretKey) > 4 {
		truncatedSecretKey = currentSecretKey[len(currentSecretKey)-4 : len(currentSecretKey)]
	}

	fmt.Printf("Ion Channel Api Key [%s]: ", truncatedSecretKey)
	secretKey, _ := terminal.ReadPassword(int(os.Stdin.Fd()))

	Debugf("All you keys are belong to us! (%s)", secretKey)

	if len(secretKey) != 0 {
		saveCredentials(string(secretKey))
	}
}

func LoadCredential() string {
	credential := os.Getenv(CREDENTIALS_ENVIRONMENT_VARIABLE)
	if credential == "" {
		Debugln("Credential env var not found looking in file")
		exists, _ := PathExists(ION_HOME)
		if exists {
			bytes, _ := ReadBytesFromFile(CREDENTIALS_FILE)
			credentials := make(map[string]string)
			yaml.Unmarshal([]byte(bytes), &credentials)
			return credentials[CREDENTIALS_KEY_FIELD]
		} else {
			MkdirAll(ION_HOME, 0775)
			return ""
		}
	} else {
		Debugln("Credential env var found")
		return credential
	}
}

func saveCredentials(secretKey string) {
	credentials := make(map[string]string)
	credentials[CREDENTIALS_KEY_FIELD] = secretKey
	yamlCredentials, _ := yaml.Marshal(&credentials)
	WriteLinesToFile(CREDENTIALS_FILE, []string{string(yamlCredentials)}, 0600)
}
