package ionconnect

import (
  "github.com/ion-channel/ion-connect/Godeps/_workspace/src/github.com/codegangsta/cli"
  "strings"
  "reflect"
)

type PostParams struct {
    Name      string   `json:"name,omitempty"`
    Url       string   `json:"url,omitempty"`
    Type      string   `json:"type,omitempty"`
    Checksum  string   `json:"checksum,omitempty"`
}

type GetParams struct {
    Name      string   `url:"name,omitempty"`
    Url       string   `url:"url,omitempty"`
    Type      string   `url:"type,omitempty"`
    Checksum  string   `url:"checksum,omitempty"`
}

func (params GetParams) Generate(context *cli.Context) GetParams {
  flags := context.Command.Flags
  for index := range flags {
    flag := flags[index]
    if flag, ok := flag.(cli.StringFlag); ok {
      reflect.ValueOf(&params).Elem().FieldByName(strings.Title(cli.StringFlag(flag).Name)).SetString(context.String(cli.StringFlag(flag).Name))
    }
  }
  return params
}

func (params PostParams) Generate(context *cli.Context) PostParams {
  flags := context.Command.Flags
  for index := range flags {
    flag := flags[index]
    if flag, ok := flag.(cli.StringFlag); ok {
      reflect.ValueOf(&params).Elem().FieldByName(strings.Title(cli.StringFlag(flag).Name)).SetString(context.String(cli.StringFlag(flag).Name))
    }
  }
  return params
}