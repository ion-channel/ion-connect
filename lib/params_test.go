// params_test.go
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

var _ = Describe("Params", func() {
	var (
		context *cli.Context
		config  Command
		set     *flag.FlagSet
	)

	BeforeEach(func() {
		Debug = true
		set = flag.NewFlagSet("set", 0)
		command := cli.Command{Name: "analyze-project"}
		context = cli.NewContext(nil, set, nil)
		context.Command = command

		config, _ = GetConfig().FindSubCommandConfig("scanner", "analyze-project")
	})

	Context("When generating Post Params", func() {
		It("should populate the fields from the context flags", func() {
			args := []string{"ernie"}
			params := PostParams{}.Generate(args, config.Args)
			Expect(params).To(Equal(PostParams{BuildNumber: "ernie"}))
			Expect(params.String()).To(Equal("Project=, Url=, Type=, Checksum=, Id=, Text=, Version=, File="))
		})

		It("should populate the fields from json data", func() {
			config, _ = GetConfig().FindSubCommandConfig("test", "test-json")
			args := []string{"{\"key\":\"value\"}", "[\"key\",\"value\"]"}
			params := PostParams{}.Generate(args, config.Args)
			m := make(map[string]interface{})
			m["key"] = "value"

			a := make([]interface{}, 2)
			a[0] = "key"
			a[1] = "value"

			Expect(params).To(Equal(PostParams{Results: m, Rules: a}))
		})
	})

	Context("When generating Get Params", func() {
		It("should populate the fields from the context flags", func() {
			args := []string{"ernie"}
			options := map[string]string{
				"limit":  "22",
				"offset": "105",
			}
			params := GetParams{}.Generate(args, config.Args).UpdateFromMap(options)
			Expect(params).To(Equal(GetParams{BuildNumber: "ernie", Limit: "22", Offset: "105"}))
			Expect(params.String()).To(Equal("Project=, Url=, Type=, Checksum=, Id=, Text=, Version=, Limit=22, Offset=105"))
		})
	})

	AfterEach(func() {
		Debug = false
	})
})
