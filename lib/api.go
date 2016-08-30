// api.go
//
// Copyright [2016] [Selection Pressure]
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ionconnect

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/dghubble/sling"
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
	args, options, err := api.validateFlags(subcommandConfig, ctx)
	if err != nil {
		fmt.Println(err.Error())
		cli.ShowCommandHelp(ctx, ctx.Command.Name)
		Exit(1)
	}

	err = api.validateArgs(args, ctx)
	if err != nil {
		fmt.Println(err.Error())
		cli.ShowCommandHelp(ctx, ctx.Command.Name)
		panic(Exit(1))
	}

	var response http.Response
	var body map[string]interface{}
	if subcommandConfig.Method == "file" {
		response, body = api.postFile(command, subcommand, ctx, args, options)
	} else {
		response, body = api.sendRequest(command, subcommand, ctx, args, options, subcommandConfig.Method)
	}
	fmt.Println(api.processResponse(response, body))
}

func (api Api) sendRequest(command string, subcommand string, context *cli.Context, args Args, options map[string]string, httpMethod string) (http.Response, map[string]interface{}) {
	client := sling.New()
	if Insecure {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		httpClient := &http.Client{Transport: transport}
		client.Client(httpClient)
	}
	var url string

	if httpMethod == "post" {
		body := PostParams{}.Generate(context.Args(), args)
		client.Post(api.Config.LoadEndpoint())
		params := GetParams{}.UpdateFromMap(options)
		client.QueryStruct(&params)
		Debugf("Sending body %s", &body)

		client.BodyJSON(&body)
		Debugf("Sending params %s", params)
	} else if httpMethod == "get" {
		params := GetParams{}.Generate(context.Args(), args).UpdateFromMap(options)
		client.Get(api.Config.LoadEndpoint())
		client.QueryStruct(&params)
		Debugf("Sending params %b", params)
	}

	url, err := api.Config.ProcessUrlFromConfig(command, subcommand, GetParams{}.Generate(context.Args(), args))
	if err != nil {
		log.Fatal(err.Error())
	}
	Debugf("Processing url: %s", url)

	Debugf("Done")

	client.Path(fmt.Sprintf("%s%s", api.Config.Version, url))
	Debugf("Sending request to %s%s", api.Config.LoadEndpoint(), fmt.Sprintf("%s%s", api.Config.Version, url))
	client.Add(api.Config.Token, LoadCredential())

	responseBody := make(map[string]interface{})
	response, responseErr := client.Receive(&responseBody, &responseBody)
	if responseErr != nil {
		Debugf("Failure occurred during request %s", responseErr.Error())
		Exit(1)
	}
	Debugf("Response received with status %s, %v", response.Status, responseBody)
	return *response, responseBody
}

func (api Api) processResponse(response http.Response, body map[string]interface{}) string {
	if response.StatusCode == 401 || response.StatusCode == 403 {
		fmt.Println("Unauthorized, make sure you run 'ion-connect configure' and set your Api Token")
		Exit(1)
		return body["message"].(string)
	} else if response.StatusCode == 400 || response.StatusCode == 404 || response.StatusCode == 422 {
		fmt.Println(body["message"])
		Exit(1)
		return body["message"].(string)
	}

	jsonBytes, err := json.MarshalIndent(body["data"], "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}

	return strings.Replace(strings.Replace(strings.Replace(string(jsonBytes), "\\u003e", ">", -1), "\\u003c", "<", -1), "\\u0026", "&", -1)
}

func (api Api) validateArgs(args Args, ctx *cli.Context) error {
	if args.GetRequiredArgsCount() > len(ctx.Args()) {
		Debugf("Missing required argument ")
		return errors.New(fmt.Sprintf("Missing required argument"))
	}

	return nil
}

func (api Api) validateFlags(commandConfig Command, ctx *cli.Context) ([]Arg, map[string]string, error) {
	args := []Arg{}
	params := make(map[string]string)
	for _, flag := range commandConfig.Flags {
		Debugf("Found option %s that is required (%b) with value %s", flag.Name, flag.Required, flag.Value)
		if !ctx.IsSet(flag.Name) && flag.Required {
			if flag.Value == "" {
				Debugf("Missing required option %s", flag.Name)
				return args, params, errors.New(fmt.Sprintf("Missing required option %s", flag.Name))
			} else {
				Debugf("Found default value for flag %s: %s", flag.Name, flag.Value)
				params[flag.Name] = flag.Value
			}
		} else if len(flag.Args) > 0 && ctx.IsSet(flag.Name) {
			Debugf("Getting args for flag %s", flag.Name)
			args = flag.Args
		} else if ctx.IsSet(flag.Name) {
			Debugf("Getting values for flag %s: %s", flag.Name, ctx.String(flag.Name))
			params[flag.Name] = ctx.String(flag.Name)
		}
	}

	if len(args) == 0 {
		args = commandConfig.Args
	}

	Debugf("Found args %v", args)
	return args, params, nil
}

func (api Api) postFile(command string, subcommand string, context *cli.Context, args Args, options map[string]string) (http.Response, map[string]interface{}) {
	params := GetParams{}.Generate(context.Args(), args).UpdateFromMap(options)
	bodyParams := PostParams{}.Generate(context.Args(), args)
	Debugf("Sending params %b", params)

	Debugf("Processing url")
	url, err := api.Config.ProcessUrlFromConfig(command, subcommand, GetParams{}.Generate(context.Args(), args))
	if err != nil {
		log.Fatal(err.Error())
	}

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	Debugf("Sending file %s", bodyParams.File)

	fileWriter, err := bodyWriter.CreateFormFile("file", bodyParams.File)
	if err != nil {
		fmt.Println("error writing to buffer")
	}

	fh, err := os.Open(bodyParams.File)
	if err != nil {
		fmt.Println(err.Error())
		Exit(1)
	}

	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		log.Fatal(err.Error())
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	url = fmt.Sprintf("%s%s%s", api.Config.LoadEndpoint(), api.Config.Version, url)
	Debugf("Sending request to %s", url)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, bodyBuf)
	req.Header.Set(api.Config.Token, LoadCredential())
	req.Header.Set("Content-Type", contentType)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
	}
	var jsonResponse map[string]interface{}
	err = json.Unmarshal([]byte(resp_body), &jsonResponse)
	if err != nil {
		panic(fmt.Sprintf("Error parsing json from %s - %s", resp_body, err.Error()))
	}
	return *resp, jsonResponse
}
