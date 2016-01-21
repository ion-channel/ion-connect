// params.go
//
// Copyright (C) 2015 Selection Pressure LLC
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package ionconnect

import (
  // "github.com/ion-channel/ion-connect/Godeps/_workspace/src/github.com/codegangsta/cli"
  "strings"
  "reflect"
)

type Params interface {
}

type PostParams struct {
    Project   string   `json:"project,omitempty"`
    Url       string   `json:"url,omitempty"`
    Type      string   `json:"type,omitempty"`
    Checksum  string   `json:"checksum,omitempty"`
    Id        string   `json:"id,omitempty"`
    Text      string   `json:"text,omitempty"`
}

type GetParams struct {
    Project   string   `url:"project,omitempty"`
    Url       string   `url:"url,omitempty"`
    Type      string   `url:"type,omitempty"`
    Checksum  string   `url:"checksum,omitempty"`
    Id        string   `url:"id,omitempty"`
    Text      string   `url:"text,omitempty"`
}

func (params GetParams) Generate(args []string, argConfigs []Arg) GetParams {
  for index, arg := range args {
    reflect.ValueOf(&params).Elem().FieldByName(strings.Title(argConfigs[index].Name)).SetString(arg)
  }
  return params
}

func (params PostParams) Generate(args []string, argConfigs []Arg) PostParams {
  for index, arg := range args {
    reflect.ValueOf(&params).Elem().FieldByName(strings.Title(argConfigs[index].Name)).SetString(arg)
  }
  return params
}
