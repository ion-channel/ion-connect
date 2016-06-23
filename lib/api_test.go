// api_test.go
//
// Copyright (C) 2015 Selection Pressure LLC
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package ionconnect

import (
	"flag"
	"github.com/codegangsta/cli"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Api", func() {
	var ()

	BeforeEach(func() {
		Debug = true
	})

	Context("If we have a command request", func() {
		It("should fail if a required argument is missing", func() {
			Test = true
			api := Api{Config: GetConfig()}
			subCommand := cli.Command{Name: "get-languages"}
			command := cli.Command{Name: "metadata", Subcommands: []cli.Command{subCommand}}
			context := cli.NewContext(nil, nil, nil)
			context.Command = command
			Expect(func() { api.HandleCommand(context) }).To(Panic())
		})

		It("should fail if an required flag is missing", func() {
			Test = true
			api := Api{Config: GetConfig()}
			set := flag.NewFlagSet("set", 0)
			set.Parse([]string{"test", "test1"})
			Expect(set.Args()).To(Equal([]string{"test", "test1"}))

			subCommand := cli.Command{Name: "test1"}
			command := cli.Command{Name: "test", Subcommands: []cli.Command{subCommand}}
			context := cli.NewContext(nil, set, nil)
			context.Command = command
			Expect(func() { api.HandleCommand(context) }).To(Panic())
		})

		It("should send the request if everything is there", func() {
			Test = true
			api := Api{Config: GetConfig()}
			set := flag.NewFlagSet("set", 0)

			subCommand := cli.Command{Name: "test1"}
			command := cli.Command{Name: "test", Subcommands: []cli.Command{subCommand}}
			context := cli.NewContext(nil, set, nil)
			context.Command = command

			config, _ := GetConfig().FindSubCommandConfig("test", "test1")
			response, body := api.sendRequest("test", "test1", context, config.Args, make(map[string]string), "get")
			Expect(response.Status).To(Equal("404 Not Found"))
			Expect(body["message"]).To(Equal("API not found with these values"))
		})

		It("should process the response body", func() {
			Test = true
			api := Api{Config: GetConfig()}
			set := flag.NewFlagSet("set", 0)

			subCommand := cli.Command{Name: "test1"}
			command := cli.Command{Name: "test", Subcommands: []cli.Command{subCommand}}
			context := cli.NewContext(nil, set, nil)
			context.Command = command

			config, _ := GetConfig().FindSubCommandConfig("test", "test1")
			response, body := api.sendRequest("test", "test1", context, config.Args, make(map[string]string), "get")
			Expect(api.processResponse(response, body)).To(Equal("API not found with these values"))
		})
	})

	AfterEach(func() {
		Debug = false
	})
})
