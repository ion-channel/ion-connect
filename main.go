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
)

func getApp() *cli.App {
	app := cli.NewApp()
	app.Name = "ion-connect"
	app.Usage = "Interact with Ion Channel"
	app.Version = "0.8.1"
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

func getSubcommands(subcommands []ionconnect.Command, api ionconnect.Api) []cli.Command {
	subs := make([]cli.Command, len(subcommands))
	for commandIndex, subcommand := range subcommands {

		flags := getFlags(subcommand.Flags)

		subs[commandIndex] = cli.Command{
			Name:      subcommand.Name,
			Usage:     subcommand.Usage,
			Action:    api.HandleCommand,
			ArgsUsage: subcommand.GetArgsUsage(),
			Flags:     flags,
		}
	}
	return subs
}

func getCommands(api ionconnect.Api) []cli.Command {
	commands := make([]cli.Command, len(api.Config.Commands)+1)

	for index, configCommand := range api.Config.Commands {
		subcommands := getSubcommands(configCommand.Subcommands, api)
		commands[index] = cli.Command{
			Name:        configCommand.Name,
			Usage:       configCommand.Usage,
			Action:      api.Noop,
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

	app.Commands = getCommands(api)

	defer deferer()
	app.Run(os.Args)
}
