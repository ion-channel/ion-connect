// util_test.go
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
	"gopkg.in/mattes/go-expand-tilde.v1"
)

var _ = Describe("Util", func() {
	var ()

	BeforeEach(func() {
		Debug = true
	})

	Context("When generating the MD5 hash", func() {
		It("should return empty if the file doesnt exist", func() {
			md5, err := ComputeMd5("/aint/real")
			Expect(err).NotTo(Equal(nil))
			Expect(md5).To(Equal(""))
		})
		It("should return MD5 if the file exists", func() {
			md5, err := ComputeMd5("../test/analysisstatus.json")
			Expect(err).To(BeNil())
			Expect(md5).To(Equal("dc8f02f8d1bd65675a609b6a2b93a943"))
		})
	})

	Context("When the debug flag is set", func() {
		It("should write out debug statements", func() {
			Expect(func() { Debugln("testing") }).ShouldNot(Panic())
			Expect(func() { Debugf("testing %s", "f") }).ShouldNot(Panic())
		})
	})

	Context("When a checking for a file or folder", func() {
		It("return false if it doesn't exist", func() {
			Expect(PathExists("/aint/real")).To(BeFalse())
		})
		It("return true if it exists", func() {
			path, _ := tilde.Expand("~")
			Expect(PathExists(path)).To(BeTrue())
		})
	})

	Context("When Running a test", func() {
		It("should not exit", func() {
			Test = true
			Expect(PathExists("/aint/real")).To(BeFalse())
		})
		It("return true if it exists", func() {
			path, _ := tilde.Expand("~")
			Expect(PathExists(path)).To(BeTrue())
		})
	})

	Context("When a param is a file type", func() {
		It("should upload a file and change the param to a url", func() {
			url := ConvertFileToUrl("file://./util.go")
			Expect(url).To(Equal("https://s3.amazonaws.com/files.ionchannel.io/files/upload/util.go"))
		})
	})

	AfterEach(func() {
		Debug = false
	})
})
