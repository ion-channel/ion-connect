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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Params_Files", func() {
	var (
		config  Command
	)

	BeforeEach(func() {
		Debug = true
		config, _ = GetConfig().FindSubCommandConfig("test", "url")
	})

	Context("When scanning a file", func() {
		It("should populate the checksum based on the file automatically, and populate the checksum param", func() {
			args := []string{"test", "file://../test/analysisstatus.json"}
			Debugf("Params_Files_Test post args : %s", args)
			params := PostParams{}.Generate(args, config.Args)
			Debugf("Params_Files_Test : %s", params)
			Expect(params.String()).To(Equal("List=[], Project=, Url=https://s3.amazonaws.com/files.ionchannel.io/files/upload/analysisstatus.json, Type=, Checksum=dc8f02f8d1bd65675a609b6a2b93a943, Id=, Text=, Version=, File=, Username=, Password="))
		})
	})

	AfterEach(func() {
		Debug = false
	})
})
