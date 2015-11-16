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
  app.Usage = "Control AWS profiles"
  app.Version = "0.1"

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

  app.Commands = []cli.Command{
    // {
    //   Name:    "describe-profiles",
    //   Aliases: []string{"d"},
    //   Usage:   `Describes the list of AWS profile`,
    //   Before:  ionconnect.BeforeDescribeProfiles,
    //   Action:  ionconnect.DescribeProfiles,
    // },
    // {
    //   Name:    "describe-active-profile",
    //   Aliases: []string{"dap"},
    //   Usage:   `Describes the currently active AWS profile`,
    //   Action:  ionconnect.DescribeActiveProfile,
    // },
    // {
    //   Name:    "activate-profile",
    //   Aliases: []string{"ap"},
    //   Usage:   `Sets the currently active profile`,
    //   Before:  ionconnect.BeforeActivateProfile,
    //   Action:  ionconnect.ActivateProfile,
    //   Flags: []cli.Flag{
    //     cli.StringFlag{
    //       Name:  "profile",
    //       Usage: "name of the profile to activate",
    //       Value: "profile-name",
    //     },
    //   },
    // },
    // {
    //   Name:    "deactive-profile",
    //   Aliases: []string{"dp"},
    //   Usage:   `Deactivate the currently active AWS profile`,
    //   Before:  ionconnect.BeforeDeactivateProfile,
    //   Action:  ionconnect.DeactivateProfile,
    // },
  }

  defer func() {
    if r := recover(); r != nil {
      fmt.Println(r)
    }
  }()
  app.Run(os.Args)
}
