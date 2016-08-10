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

var _ = Describe("Params_Files", func() {
	var (
		context *cli.Context
		config  Command
		set     *flag.FlagSet
	)

	BeforeEach(func() {
		Debug = true
		set = flag.NewFlagSet("set", 0)
		command := cli.Command{Name: "scan-artifact-url"}
		context = cli.NewContext(nil, set, nil)
		context.Command = command

		config, _ = GetConfig().FindSubCommandConfig("scanner", "scan-artifact-url")
	})

	Context("When scanning a file", func() {
		It("should populate the checksum based on the file automatically, and populate the checksum param", func() {
			args := []string{"fart", "file://../test/analysisstatus.json"}
			//args := []string{"{\"project\":\"fart\"}", "[\"url\",\"file://../test/analysisstatus.json\"]"}
			Debugf("Params_Files_Test post args : %s", args)
			params := PostParams{}.Generate(args, config.Args)
			Debugf("Params_Files_Test : %s", params)
			Expect(params.String()).To(Equal("Project=fart, Url=https://s3.amazonaws.com/files.ionchannel.io/files/upload/analysisstatus.json, Type=, Checksum=29803548a0cf1281078bb9d88621ddb8, Id=, Text=, Version=, File="))
		})
	})

	AfterEach(func() {
		Debug = false
	})
})
