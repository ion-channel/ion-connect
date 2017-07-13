// main.go
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

package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/ion-channel/ion-connect/lib"
	"github.com/sirupsen/logrus"
)

func getApp() *cli.App {
	app := cli.NewApp()
	app.Name = "ion-connect"
	app.Usage = "Interact with Ion Channel"
	app.Version = "0.10.2"
	return app
}

func setFlags(app *cli.App) *cli.App {
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "display debug logging",
		},
		cli.BoolFlag{
			Name:  "insecure",
			Usage: "allow for insecure https connections",
		},
	}

	return app
}

func before(c *cli.Context) error {
	ionconnect.Logger.Out = os.Stdout
	ionconnect.Logger.Level = logrus.DebugLevel
	ionconnect.Logger.Formatter = &logrus.TextFormatter{}

	if c.Bool("debug") {
		ionconnect.Debug = true
		ionconnect.Debugln("Turning debug on.")
	}

	if c.Bool("insecure") {
		ionconnect.Insecure = true
		ionconnect.Debugln("Turning insecure on.")
	}

	return nil
}

func getFlags(flagConfigs []ionconnect.Flag) []cli.Flag {
	flags := make([]cli.Flag, len(flagConfigs))
	for flagIndex, flag := range flagConfigs {
		switch flag.Type {
		case "bool":
			flags[flagIndex] = cli.BoolFlag{
				Name:  flag.Name,
				Usage: flag.Usage,
			}
		case "string":
			flags[flagIndex] = cli.StringFlag{
				Name:  flag.Name,
				Value: flag.Value,
				Usage: flag.Usage,
			}
		}
	}
	return flags
}

func getSubcommands(subcommands []ionconnect.Command, handler interface{}) []cli.Command {
	subs := make([]cli.Command, len(subcommands))
	for commandIndex, subcommand := range subcommands {

		flags := getFlags(subcommand.Flags)

		subs[commandIndex] = cli.Command{
			Name:      subcommand.Name,
			Usage:     subcommand.Usage,
			Action:    handler,
			ArgsUsage: subcommand.GetArgsUsage(),
			Flags:     flags,
		}
	}
	return subs
}

func getCommands(configCommands []ionconnect.Command, noop interface{}, handler interface{}) []cli.Command {
	commands := make([]cli.Command, len(configCommands)+1)

	for index, configCommand := range configCommands {
		subcommands := getSubcommands(configCommand.Subcommands, handler)
		commands[index] = cli.Command{
			Name:        configCommand.Name,
			Usage:       configCommand.Usage,
			Action:      noop,
			Subcommands: subcommands,
		}
	}

	commands[len(commands)-1] = cli.Command{
		Name:   "configure",
		Usage:  "setup the Ion Channel secret key for later use",
		Action: ionconnect.HandleConfigure,
	}
	return commands
}

func deferer() {
	if r := recover(); r != nil {
		fmt.Println(r)
	}
}

func main() {
	app := getApp()

	var api = ionconnect.Api{ionconnect.GetConfig()}

	app = setFlags(app)

	app.Before = before

	app.Commands = getCommands(api.Config.Commands, api.Noop, api.HandleCommand)

	defer deferer()
	app.Run(os.Args)
}
