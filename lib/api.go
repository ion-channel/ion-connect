package ionconnect

import (
  "github.com/ion-channel/ion-connect/Godeps/_workspace/src/github.com/codegangsta/cli"
  "github.com/ion-channel/ion-connect/Godeps/_workspace/src/github.com/dghubble/sling"
  "strings"
  "log"
  "errors"
  "fmt"
)

type Api struct {
  Config Config
}

func (api Api) HandleCommand(ctx *cli.Context) {
  Debugf("Performing command %s", ctx.Command.FullName())


  command := strings.Split(ctx.Command.FullName(), " ")[0]
  subcommand := strings.Split(ctx.Command.FullName(), " ")[1]
  subcommandConfig, err := api.Config.findSubCommandConfig(command, subcommand)
  if err != nil {
    log.Fatalf("Command configuration missing for %s %s", command, subcommand)
  }
  err = api.validateFlags(subcommandConfig, ctx)
  if err != nil {
    fmt.Println(err.Error())
    cli.ShowCommandHelp(ctx, ctx.Command.Name)
  }

  client := sling.New().Base(api.Config.Endpoint)
  client.Path(command).Path(subcommand)
  client.Add(api.Config.Token, "token")
  client.Request()
}

func (api Api) Noop(ctx *cli.Context) {
  cli.ShowCommandHelp(ctx, ctx.Command.Name)
  Debugln("Noop")
}

func (api Api) validateFlags(commandConfig Command, ctx *cli.Context) error {
  for index := range commandConfig.Flags {
    if ctx.String(commandConfig.Flags[index].Name) == "" {
      Debugf("Missing required option %s", commandConfig.Flags[index].Name)
      return errors.New(fmt.Sprintf("Missing required option %s", commandConfig.Flags[index].Name))
    }
  }

  return nil
}
