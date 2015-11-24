// api.go
//
// Copyright (C) 2015 Selection Pressure LLC
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package ionconnect

import (
  "github.com/ion-channel/ion-connect/Godeps/_workspace/src/github.com/codegangsta/cli"
  "github.com/ion-channel/ion-connect/Godeps/_workspace/src/github.com/dghubble/sling"
  "encoding/json"
  "net/http"
  "strings"
  "log"
  "errors"
  "fmt"
  "os"
)

type Api struct {
  Config Config
}

func (api Api) Noop(ctx *cli.Context) {
  cli.ShowCommandHelp(ctx, ctx.Command.Name)
  Debugln("Noop")
}

func (api Api) HandleCommand(ctx *cli.Context) {
  Debugf("Performing command %s", ctx.Command.FullName())
  command := strings.Split(ctx.Command.FullName(), " ")[0]
  subcommand := strings.Split(ctx.Command.FullName(), " ")[1]
  subcommandConfig, err := api.Config.FindSubCommandConfig(command, subcommand)
  if err != nil {
    log.Fatalf("Command configuration missing for %s %s", command, subcommand)
  }
  err = api.validateFlags(subcommandConfig, ctx)
  if err != nil {
    fmt.Println(err.Error())
    cli.ShowCommandHelp(ctx, ctx.Command.Name)
    os.Exit(1)
  }

  response, body := api.sendRequest(command, subcommand, ctx, subcommandConfig.Post)

  fmt.Println(api.processResponse(response, body))
}

func (api Api) sendRequest(command string, subcommand string, context *cli.Context, should_post bool) (http.Response, map[string]interface{}) {
  client := sling.New()
  var url string

  if should_post {
    params := PostParams{}.Generate(context)
    client.Post(api.Config.Endpoint)
    client.BodyJSON(&params)
  } else {
    params := GetParams{}.Generate(context)
    client.Get(api.Config.Endpoint)
    client.QueryStruct(&params)
  }

  url, err := api.Config.ProcessUrlFromConfig(command, subcommand, GetParams{}.Generate(context))
  if err != nil {
    log.Fatal(err.Error())
  }

  client.Path(fmt.Sprintf("%s%s", api.Config.Version, url))
  client.Add(api.Config.Token, LoadCredential())
  body := make(map[string]interface{})
  response, responseErr := client.Receive(&body, &body)
  if responseErr != nil {
    fmt.Println(responseErr.Error())
    os.Exit(1)
  }
  Debugf("Response received with status %s, %v", response.Status, body)
  return *response, body
}



func (api Api) processResponse(response http.Response, body map[string]interface{}) string {
  if response.StatusCode == 401 {
    fmt.Println("Unauthorized, make sure you run 'ion-connect configure' and set your Api Token")
    os.Exit(1)
  }

  if response.StatusCode == 400 {
    fmt.Println(body["message"])
    os.Exit(1)
  }

  delete(body, "links")
  jsonBytes, err := json.Marshal(body)
  if err != nil {
    log.Fatalf(err.Error())
  }

  return string(jsonBytes)
}

func (api Api) validateFlags(commandConfig Command, ctx *cli.Context) error {
  for _, flag := range commandConfig.Flags {
    if ctx.String(flag.Name) == "" {
      Debugf("Missing required option %s", flag.Name)
      return errors.New(fmt.Sprintf("Missing required option %s", flag.Name))
    }
  }

  return nil
}
