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
    os.Exit(1)
  }

  response, body, errorBody := api.sendRequest(command, subcommand, ctx, subcommandConfig.Write)

  fmt.Println(api.processResponse(response, body, errorBody))
}

func (api Api) sendRequest(command string, subcommand string, context *cli.Context, should_post bool) (http.Response, map[string]interface{}, map[string]interface{}) {
  client := sling.New()

  if should_post {
    params := PostParams{}.Generate(context)
    client.Post(api.Config.Endpoint)
    client.BodyJSON(&params)
  } else {
    params := GetParams{}.Generate(context)
    client.Get(api.Config.Endpoint)
    client.QueryStruct(&params)
  }
  client.Path(fmt.Sprintf("%s/%s/%s", api.Config.Version, command, subcommand))
  client.Add(api.Config.Token, LoadCredential())

  body := make(map[string]interface{})
  errorBody := make(map[string]interface{})
  response, responseErr := client.Receive(&body, &errorBody)
  if responseErr != nil {
    log.Fatalf(responseErr.Error())
  }
  Debugf("Response received with status %s, %v", response.Status, errorBody)
  return *response, body, errorBody
}

func (api Api) processResponse(response http.Response, body map[string]interface{}, errorBody map[string]interface{}) string {
  if response.StatusCode == 401 {
    fmt.Println("Unauthorized, make sure you run 'ion-connect configure' and set your Api Token")
    os.Exit(1)
  }

  if response.StatusCode == 400 {
    fmt.Println(errorBody["message"])
    os.Exit(1)
  }

  delete(body, "links")
  jsonBytes, err := json.Marshal(body)
  if err != nil {
    log.Fatalf(err.Error())
  }

  return string(jsonBytes)
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
