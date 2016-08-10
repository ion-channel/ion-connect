// util_test.go
//
// Copyright (C) 2015 Selection Pressure LLC
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package ionconnect

import (
	"strings"

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
			Expect(strings.ToUpper(md5)).To(Equal("29803548A0CF1281078BB9D88621DDB8"))
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
