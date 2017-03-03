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

	Context("when flags for a command are configured", func() {
		It("should generate command line flags for a string", func() {
			expectedFlag := cli.StringFlag{
				Name:        "someflag",
				Value:       "somevalue",
				Usage:       "someusage",
				EnvVar:      "",
				Destination: nil,
				Hidden:      false,
			}

			flag := ionconnect.Flag{
				Name:     "someflag",
				Value:    "somevalue",
				Usage:    "someusage",
				Type:     "string",
				Required: true,
			}
			flags := getFlags([]ionconnect.Flag{flag})
			Expect(flags[0]).To(Equal(expectedFlag))
		})

		It("should generate command line flags for a boolean", func() {
			expectedFlag := cli.BoolFlag{
				Name:        "someflag",
				Usage:       "someusage",
				EnvVar:      "",
				Destination: nil,
				Hidden:      false,
			}

			flag := ionconnect.Flag{
				Name:     "someflag",
				Value:    "",
				Usage:    "someusage",
				Type:     "bool",
				Required: true,
			}
			flags := getFlags([]ionconnect.Flag{flag})
			Expect(flags[0]).To(Equal(expectedFlag))
		})
	})

	Context("when a subcommand is configured", func() {
		It("should generate a cli subcommand with flags", func() {
			expectedCommand := cli.Command{
				Name:        "somecommand",
				ShortName:   "",
				Aliases:     nil,
				Usage:       "someusage",
				UsageText:   "",
				Description: "",
				ArgsUsage:   "",
				Subcommands: nil,
				Flags:       []cli.Flag{cli.StringFlag{Name: "someflag", Usage: "someusage"}},
			}

			flag := ionconnect.Flag{
				Name:     "someflag",
				Value:    "",
				Usage:    "someusage",
				Type:     "string",
				Required: true,
			}

			command := ionconnect.Command{
				Name:   "somecommand",
				Usage:  "someusage",
				Method: "POST",
				Url:    "/some/url",
				Flags:  []ionconnect.Flag{flag},
			}

			Expect(getSubcommands([]ionconnect.Command{command}, nil)).To(Equal([]cli.Command{expectedCommand}))
		})

		It("should generate a cli subcommand without flags", func() {
			expectedCommand := cli.Command{
				Name:        "somecommand",
				ShortName:   "",
				Aliases:     nil,
				Usage:       "someusage",
				UsageText:   "",
				Description: "",
				ArgsUsage:   "",
				Subcommands: nil,
				Flags:       []cli.Flag{},
			}

			command := ionconnect.Command{
				Name:   "somecommand",
				Usage:  "someusage",
				Method: "POST",
				Url:    "/some/url",
			}

			Expect(getSubcommands([]ionconnect.Command{command}, nil)).To(Equal([]cli.Command{expectedCommand}))
		})
	})


	Context("when a command is configured", func() {
    It("should generate a cli command without flags", func(){
      expectedCommand := cli.Command{
        Name:        "somecommand",
        ShortName:   "",
        Aliases:     nil,
        Usage:       "someusage",
        UsageText:   "",
        Description: "",
        ArgsUsage:   "",
        Subcommands: []cli.Command{},
        Flags:       nil,
        Action: nil,
      }

      command := ionconnect.Command{
        Name:   "somecommand",
        Usage:  "someusage",
        Method: "POST",
        Url:    "/some/url",
      }
      Expect(getCommands([]ionconnect.Command{command}, nil, nil)[0]).To(Equal(expectedCommand))
    })
  })
	AfterEach(func() {
	})
})
