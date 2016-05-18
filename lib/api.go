// api.go
//
// Copyright (C) 2015 Selection Pressure LLC
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package ionconnect

import (
	"encoding/json"
	"github.com/ion-channel/ion-connect/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/ion-channel/ion-connect/Godeps/_workspace/src/github.com/dghubble/sling"
	"net/http"
	// "strings"
	"errors"
	"fmt"
	"log"
	"strings"

	"crypto/tls"
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
	response, body := api.sendRequest(command, subcommand, ctx, args, options, subcommandConfig.Post)

	fmt.Println(api.processResponse(response, body))
}

func (api Api) sendRequest(command string, subcommand string, context *cli.Context, args Args, options map[string]string, shouldPost bool) (http.Response, map[string]interface{}) {
	client := sling.New()
	if Insecure {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		httpClient := &http.Client{Transport: transport}
		client.Client(httpClient)
	}
	var url string

	if shouldPost {
		body := PostParams{}.Generate(context.Args(), args)
		client.Post(api.Config.LoadEndpoint())
		params := GetParams{}.UpdateFromMap(options)
		client.QueryStruct(&params)
		client.BodyJSON(&body)
		Debugf("Sending body %v", body)
		Debugf("Sending params %v", params)
	} else {
		params := GetParams{}.Generate(context.Args(), args).UpdateFromMap(options)
		client.Get(api.Config.LoadEndpoint())
		client.QueryStruct(&params)
		Debugf("Sending params %v", params)
	}

	Debugf("Processing url")
	url, err := api.Config.ProcessUrlFromConfig(command, subcommand, GetParams{}.Generate(context.Args(), args))
	if err != nil {
		log.Fatal(err.Error())
	}

	Debugf("Done")

	client.Path(fmt.Sprintf("%s%s", api.Config.Version, url))
	Debugf("Sending request to %s%s", api.Config.LoadEndpoint(), fmt.Sprintf("%s%s", api.Config.Version, url))
	client.Add(api.Config.Token, LoadCredential())

	body := make(map[string]interface{})
	response, responseErr := client.Receive(&body, &body)
	if responseErr != nil {
		fmt.Println(responseErr.Error())
		Exit(1)
	}
	Debugf("Response received with status %s, %v", response.Status, body)
	return *response, body
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
		if !ctx.IsSet(flag.Name) && flag.Required {
			Debugf("Missing required option %s", flag.Name)
			return args, params, errors.New(fmt.Sprintf("Missing required option %s", flag.Name))
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
