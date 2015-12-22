// main.go
//
// Copyright (C) 2015 Selection Pressure LLC
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package main

import (
	"fmt"
	"github.com/ion-channel/ion-connect/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/ion-channel/ion-connect/lib"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "ionconnect"
	app.Usage = "Interact with Ion Channel"
	app.Version = "0.2.0"

	var api = ionconnect.Api{ionconnect.GetConfig()}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "display debug logging",
		},
	}
	app.Before = func(c *cli.Context) error {
		if c.Bool("debug") {
			ionconnect.Debug = true
			ionconnect.Debugln("Turning debug on.")
		}

		return nil
	}

	commands := make([]cli.Command, len(api.Config.Commands)+1)

	for index, configCommand := range api.Config.Commands {
		subcommands := make([]cli.Command, len(configCommand.Subcommands))
		for jndex, subcommand := range configCommand.Subcommands {

			flags := make([]cli.Flag, len(subcommand.Flags))
			for kndex, flag := range subcommand.Flags {
				flags[kndex] = cli.StringFlag{
					Name:  flag.Name,
					Value: flag.Value,
					Usage: flag.Usage,
				}
			}

			subcommands[jndex] = cli.Command{
				Name:   subcommand.Name,
				Usage:  subcommand.Usage,
				Action: api.HandleCommand,
				Flags:  flags,
			}
		}
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

	app.Commands = commands

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	app.Run(os.Args)
}
