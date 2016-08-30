// params_test.go
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
			Expect(params.String()).To(Equal("List=[], Project=fart, Url=https://s3.amazonaws.com/files.ionchannel.io/files/upload/analysisstatus.json, Type=, Checksum=29803548a0cf1281078bb9d88621ddb8, Id=, Text=, Version=, File="))
		})
	})

	AfterEach(func() {
		Debug = false
	})
})
