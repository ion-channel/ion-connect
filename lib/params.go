// params.go
//
// Copyright (C) 2015 Selection Pressure LLC
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package ionconnect

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Params interface {
}

type PostParams struct {
	Project     string                 `json:"project,omitempty"`
	Product     string                 `json:"product,omitempty"`
	Url         string                 `json:"url,omitempty"`
	Type        string                 `json:"type,omitempty"`
	Checksum    string                 `json:"checksum,omitempty"`
	Id          string                 `json:"id,omitempty"`
	Text        string                 `json:"text,omitempty"`
	Version     string                 `json:"version,omitempty"`
	File        string                 `json:"file,omitempty"`
	ProjectId   string                 `json:"project_id,omitempty"`
	AccountId   string                 `json:"account_id,omitempty"`
	AnalysisId  string                 `json:"analysis_id,omitempty"`
	ScanId      string                 `json:"scan_id,omitempty"`
	RulesetId   string                 `json:"ruleset_id,omitempty"`
	BuildNumber string                 `json:"build_number,omitempty"`
	Status      string                 `json:"status,omitempty"`
	Results     map[string]interface{} `json:"results,omitempty"`
	ScanType    string                 `json:"scan_type,omitempty"`
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Branch      string                 `json:"branch,omitempty"`
	Source      string                 `json:"source,omitempty"`
	Active      bool                   `json:"active,omitempty"`
	Flatten     bool                   `json:"flatten,omitempty"`
	Rules       []interface{}          `json:"rules,omitempty"`
	RuleIds     []interface{}          `json:"rule_ids,omitempty"`
	ScanSet     []interface{}          `json:"data,omitempty"`
	SkipAck     bool                   `json:"skip_ack,omitempty"`
}

type GetParams struct {
	Project     string                 `url:"project,omitempty"`
	Product     string                 `url:"product,omitempty"`
	Url         string                 `url:"url,omitempty"`
	Type        string                 `url:"type,omitempty"`
	Checksum    string                 `url:"checksum,omitempty"`
	Id          string                 `url:"id,omitempty"`
	Text        string                 `url:"text,omitempty"`
	Version     string                 `url:"version,omitempty"`
	Limit       string                 `url:"limit,omitempty"`
	Offset      string                 `url:"offset,omitempty"`
	File        string                 `url:"file,omitempty"`
	ProjectId   string                 `url:"project_id,omitempty"`
	AccountId   string                 `url:"account_id,omitempty"`
	AnalysisId  string                 `url:"analysis_id,omitempty"`
	ScanId      string                 `url:"scan_id,omitempty"`
	RulesetId   string                 `url:"ruleset_id,omitempty"`
	BuildNumber string                 `url:"build_number,omitempty"`
	Status      string                 `url:"status,omitempty"`
	Results     map[string]interface{} `url:"results,omitempty"`
	ScanType    string                 `url:"scan_type,omitempty"`
	Name        string                 `url:"name,omitempty"`
	Description string                 `url:"description,omitempty"`
	Branch      string                 `url:"branch,omitempty"`
	Source      string                 `url:"source,omitempty"`
	Active      bool                   `url:"active,omitempty"`
	Flatten     bool                   `url:"flatten,omitempty"`
	Rules       []interface{}          `url:"rules,omitempty"`
	RuleIds     []interface{}          `url:"rule_ids,omitempty"`
	ScanSet     []interface{}          `url:"data,omitempty"`
	SkipAck     bool                   `url:"skip_ack,omitempty"`
}

func (params GetParams) String() string {
	return fmt.Sprintf("Project=%s, Url=%s, Type=%s, Checksum=%s, Id=%s, Text=%s, Version=%s, Limit=%s, Offset=%s", params.Project, params.Url, params.Type, params.Checksum, params.Id, params.Text, params.Version, params.Limit, params.Offset)
}

func (params GetParams) Generate(args []string, argConfigs []Arg) GetParams {
	for index, arg := range args {
		if argConfigs[index].Type != "object" && argConfigs[index].Type != "array" {
			reflect.ValueOf(&params).Elem().FieldByName(strings.Replace(strings.Title(argConfigs[index].Name), "-", "", -1)).SetString(arg)
		} else if argConfigs[index].Type == "bool" {
			boolArg, _ := strconv.ParseBool(arg)
			reflect.ValueOf(&params).Elem().FieldByName(strings.Replace(strings.Title(argConfigs[index].Name), "-", "", -1)).SetBool(boolArg)
		}
	}
	return params
}

func (params GetParams) UpdateFromMap(paramMap map[string]string) GetParams {
	for param, value := range paramMap {
		Debugf("Setting %s to %s", strings.Replace(strings.Title(param), "-", "", -1), value)
		_, intErr := strconv.ParseInt(value, 10, 64)
		boolValue, boolErr := strconv.ParseBool(value)
		if boolErr == nil && intErr != nil {
			reflect.ValueOf(&params).Elem().FieldByName(strings.Replace(strings.Title(param), "-", "", -1)).SetBool(boolValue)
		} else {
			reflect.ValueOf(&params).Elem().FieldByName(strings.Replace(strings.Title(param), "-", "", -1)).SetString(value)
		}
	}
	return params
}

func (params PostParams) String() string {
	return fmt.Sprintf("Project=%s, Url=%s, Type=%s, Checksum=%s, Id=%s, Text=%s, Version=%s, File=%s", params.Project, params.Url, params.Type, params.Checksum, params.Id, params.Text, params.Version, params.File)
}

func (params PostParams) Generate(args []string, argConfigs []Arg) PostParams {
	var md5hash string
	for index, arg := range args {
		Debugf("Index and args %d %s %v", index, arg, argConfigs)

		Debugf("PostParams Setting %s to %s", strings.Title(argConfigs[index].Name), arg)
		if argConfigs[index].Type == "object" {
			Debugln("Using object parser")
			var jsonArg map[string]interface{}
			err := json.Unmarshal([]byte(arg), &jsonArg)
			if err != nil {
				panic(fmt.Sprintf("Error parsing json from %s - %s", argConfigs[index].Name, err.Error()))
			}
			reflect.ValueOf(&params).Elem().FieldByName(strings.Replace(strings.Title(argConfigs[index].Name), "-", "", -1)).Set(reflect.ValueOf(jsonArg))
		} else if argConfigs[index].Type == "array" {
			Debugln("Using array parser")
			var jsonArray []interface{}
			err := json.Unmarshal([]byte(arg), &jsonArray)
			if err != nil {
				panic(fmt.Sprintf("Error parsing json from %s - %s", argConfigs[index].Name, err.Error()))
			}
			reflect.ValueOf(&params).Elem().FieldByName(strings.Replace(strings.Title(argConfigs[index].Name), "-", "", -1)).Set(reflect.ValueOf(jsonArray))
		} else if argConfigs[index].Type == "bool" {
			Debugf("Using bool parser for (%s) = (%s)", argConfigs[index].Name, arg)
			if arg == "" {
				Debugf("Missing arg value (%s) using default (%s)", argConfigs[index].Name, argConfigs[index].Value)
				arg = argConfigs[index].Value
			}
			boolArg, _ := strconv.ParseBool(arg)
			reflect.ValueOf(&params).Elem().FieldByName(strings.Replace(strings.Title(argConfigs[index].Name), "-", "", -1)).SetBool(boolArg)
		} else {
			if argConfigs[index].Type == "url" {
				Debugf("Handling url %s", arg)
				a, err := ComputeMd5(arg)
				md5hash = a
				if err != nil {
					fmt.Printf("Failed to generate MD5 from url %s. Make sure the file exists and permissions are correct. (%s)", arg, err)
					Exit(1)
				}
				arg = ConvertFileToUrl(arg)
			}
			Debugf("Using string parser for (%s) = (%s)", argConfigs[index].Name, arg)
			if arg == "" {
				Debugf("Missing arg value (%s) using default (%s)", argConfigs[index].Name, argConfigs[index].Value)
				arg = argConfigs[index].Value
			}
			reflect.ValueOf(&params).Elem().FieldByName(strings.Replace(strings.Title(argConfigs[index].Name), "-", "", -1)).SetString(arg)
		}

		Debugf("Finished %s", arg)
	}
	if len(md5hash) > 0 {
		params.Checksum = md5hash
	}
	return params
}
