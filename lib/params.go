// params.go
//
// Copyright (C) 2015 Selection Pressure LLC
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package ionconnect

import (
  "github.com/ion-channel/ion-connect/Godeps/_workspace/src/github.com/codegangsta/cli"
  "strings"
  "reflect"
)

type Params interface {
}

type PostParams struct {
    Name      string   `json:"name,omitempty"`
    Url       string   `json:"url,omitempty"`
    Type      string   `json:"type,omitempty"`
    Checksum  string   `json:"checksum,omitempty"`
    Scanid    string   `json:"-"`
}

type GetParams struct {
    Name      string   `url:"name,omitempty"`
    Url       string   `url:"url,omitempty"`
    Type      string   `url:"type,omitempty"`
    Checksum  string   `url:"checksum,omitempty"`
    Scanid    string   `url:"-"`
}

func (params GetParams) Generate(context *cli.Context) GetParams {
  flags := context.Command.Flags
  for _, flag := range flags {
    if flag, ok := flag.(cli.StringFlag); ok {
      reflect.ValueOf(&params).Elem().FieldByName(strings.Title(cli.StringFlag(flag).Name)).SetString(context.String(cli.StringFlag(flag).Name))
    }
  }
  return params
}

func (params PostParams) Generate(context *cli.Context) PostParams {
  flags := context.Command.Flags
  for _, flag := range flags {
    if flag, ok := flag.(cli.StringFlag); ok {
      reflect.ValueOf(&params).Elem().FieldByName(strings.Title(cli.StringFlag(flag).Name)).SetString(context.String(cli.StringFlag(flag).Name))
    }
  }
  return params
}
