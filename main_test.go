// config_test.go
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

package main

import (
	// "fmt"
	// "os"
	"flag"
	"github.com/codegangsta/cli"
	"github.com/ion-channel/ion-connect/lib"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestMain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ion Connect Main Test Suite")
}

var _ = Describe("Main", func() {
	var ()

	BeforeEach(func() {
		ionconnect.Insecure = false
		ionconnect.Debug = false
	})

	Context("when an error occurs it should be recovered and logged", func() {
		It("should log any error that is not handled", func() {
			defer deferer()
			panic("not really")
		})
	})

	Context("when an app object is requested", func() {
		It("should initialize and create an app object", func() {
			app := getApp()
			Expect(app.Name).To(Equal("ion-connect"))
			Expect(app.Usage).To(Equal("Interact with Ion Channel"))
			Expect(app.Version).NotTo(Equal(""))
		})
	})

	Context("when an app object is present", func() {
		It("should allow for flags to be attached", func() {
			app := getApp()
			app = setFlags(app)
			Expect(len(app.Flags)).To(Equal(2))
		})
	})

	Context("when an app flag is expected", func() {
		It("should be handled if it is set", func() {
			set := flag.NewFlagSet("test", 0)
			set.Bool("debug", true, "Debug logging")
			set.Bool("insecure", true, "Insecure comms")
			c := cli.NewContext(nil, set, nil)
			before(c)
			Expect(ionconnect.Debug).To(BeTrue())
			Expect(ionconnect.Insecure).To(BeTrue())
		})

		It("should not be handled if it is set to false", func() {
			set := flag.NewFlagSet("test", 0)
			set.Bool("insecure", false, "Insecure")
			set.Bool("debug", false, "Debug")
			c := cli.NewContext(nil, set, nil)
			before(c)
			Expect(ionconnect.Insecure).To(BeFalse())
			Expect(ionconnect.Debug).To(BeFalse())
		})

		It("should not be handled if flag is not set", func() {
			set := flag.NewFlagSet("test", 0)
			c := cli.NewContext(nil, set, nil)
			before(c)
			Expect(ionconnect.Insecure).To(BeFalse())
			Expect(ionconnect.Debug).To(BeFalse())
		})
	})

	AfterEach(func() {
	})
})
