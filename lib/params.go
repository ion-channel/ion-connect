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
  "fmt"
  "io/ioutil"
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
    Version   string   `json:"version,omitempty"`
    File      string   `json:"file,omitempty"`
}

type GetParams struct {
    Project   string   `url:"project,omitempty"`
    Url       string   `url:"url,omitempty"`
    Type      string   `url:"type,omitempty"`
    Checksum  string   `url:"checksum,omitempty"`
    Id        string   `url:"id,omitempty"`
    Text      string   `url:"text,omitempty"`
    Version   string   `url:"version,omitempty"`
    Limit     string   `url:"limit,omitempty"`
    Offset    string   `url:"offset,omitempty"`
    File      string   `url:"file,omitempty"`
}

func (params GetParams) String() string {
  return fmt.Sprintf("Project=%s, Url=%s, Type=%s, Checksum=%s, Id=%s, Text=%s, Version=%s, Limit=%s, Offset=%s", params.Project, params.Url, params.Type, params.Checksum, params.Id, params.Text, params.Version, params.Limit, params.Offset)
}

func (params GetParams) Generate(args []string, argConfigs []Arg) GetParams {
  for index, arg := range args {
    reflect.ValueOf(&params).Elem().FieldByName(strings.Title(argConfigs[index].Name)).SetString(arg)
  }
  return params
}

func (params GetParams) UpdateFromMap(paramMap map[string]string) GetParams {
  for param, value := range paramMap {
    Debugf("Setting %s to %s", param, value )
    reflect.ValueOf(&params).Elem().FieldByName(strings.Title(param)).SetString(value)
  }
  return params
}

func (params PostParams) String() string {
  return fmt.Sprintf("Project=%s, Url=%s, Type=%s, Checksum=%s, Id=%s, Text=%s, Version=%s, File=%s", params.Project, params.Url, params.Type, params.Checksum, params.Id, params.Text, params.Version, params.File)
}

func (params PostParams) Generate(args []string, argConfigs []Arg) PostParams {
  for index, arg := range args {
    Debugf("Setting %s to %s", argConfigs[index].Name, arg )
    if argConfigs[index].Type == "file" {
      Debugf("Reading file %s", arg)
      bytes, err := ioutil.ReadFile(arg)
      if err != nil {
        panic(err.Error())
      }
      arg = string(bytes)
    }
    reflect.ValueOf(&params).Elem().FieldByName(strings.Title(argConfigs[index].Name)).SetString(arg)
  }
  return params
}
