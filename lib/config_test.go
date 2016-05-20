// config_test.go
//
// Copyright (C) 2015 Selection Pressure LLC
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package ionconnect

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Config", func() {
	var ()

	BeforeEach(func() {
		Debug = false
		ION_HOME = fmt.Sprintf("%s/ionchannel-test", os.TempDir())
		CREDENTIALS_FILE = fmt.Sprintf("%s/credentials", ION_HOME)
	})

	Context("When the config.yaml file is loaded", func() {
		config := GetConfig()
		It("should contain an 'api' section", func() {
			Expect(len(config.Commands)).To(Equal(11))
		})
		It("should contain the api version", func() {
			Expect(config.Version).To(Equal("v1"))
		})
		It("should contain the api endpoint", func() {
			Expect(config.Endpoint).To(Equal("https://api.ionchannel.io/"))
		})
		It("should contain the api token header name", func() {
			Expect(config.Token).To(Equal("apikey"))
		})
		It("should have commands with subcommands", func() {
			Expect(len(config.Commands[0].Subcommands)).To(Equal(6))
			Expect(config.Commands[0].Subcommands[0].Name).To(Equal("scan-git"))
			Expect(config.Commands[0].Subcommands[0].Method).To(Equal("post"))
		})
		It("should have commands with subcommands and args", func() {
			Expect(len(config.Commands[0].Subcommands)).To(Equal(6))
			Expect(config.Commands[0].Subcommands[0].Name).To(Equal("scan-git"))
			Expect(config.Commands[0].Subcommands[0].Args[0].Name).To(Equal("project"))
		})
	})

	Context("When processing a url", func() {
		config := GetConfig()
		It("it should render template code", func() {
			params := GetParams{Id: "test"}
			url, err := config.ProcessUrlFromConfig("metadata", "get-locations", params)
			Expect(url).To(Equal("/metadata/getLocations"))
			Expect(err).To(BeNil())
		})
		It("it should not fail if it's just a string", func() {
			params := GetParams{Id: "test"}
			url, err := config.ProcessUrlFromConfig("metadata", "get-locations", params)
			Expect(url).To(Equal("/metadata/getLocations"))
			Expect(err).To(BeNil())
		})
	})

	Context("When we need the endpoint", func() {
		config := GetConfig()
		It("should read the endpoint from the config", func() {
			Expect(config.LoadEndpoint()).To(Equal(config.Endpoint))
		})
		It("should read the endpoint from the config", func() {
			Expect(config.LoadEndpoint()).To(Equal(config.Endpoint))
			os.Setenv(ENDPOINT_ENVIRONMENT_VARIABLE, "numbersandletters")
			Expect(config.LoadEndpoint()).To(Equal("numbersandletters"))
			os.Unsetenv(ENDPOINT_ENVIRONMENT_VARIABLE)
		})

	})

	Context("When we need creds", func() {
		It("should create the ION_HOME directory", func() {
			Expect(LoadCredential()).To(Equal(""))
			Expect(PathExists(ION_HOME)).To(BeTrue())
		})
		It("should save credentials to new file", func() {
			Expect(LoadCredential()).To(Equal(""))
			Expect(PathExists(ION_HOME)).To(BeTrue())
			Expect(PathExists(CREDENTIALS_FILE)).To(BeFalse())
			Expect(func() { saveCredentials("notarealkey") }).ShouldNot(Panic())
			Expect(PathExists(CREDENTIALS_FILE)).To(BeTrue())
			Expect(LoadCredential()).To(Equal("notarealkey"))
		})
		It("should read credentials from and environment variable", func() {
			os.Setenv(CREDENTIALS_ENVIRONMENT_VARIABLE, "numbersandletters")
			Expect(LoadCredential()).To(Equal("numbersandletters"))
			os.Unsetenv(CREDENTIALS_ENVIRONMENT_VARIABLE)
		})
	})

	Context("If we are looking for a command from the config", func() {
		config := GetConfig()
		It("should return the command config if found", func() {
			command, err := config.FindCommandConfig("metadata")
			Expect(command.Name).To(Equal("metadata"))
			Expect(err).To(BeNil())
		})
		It("should return an error if not found", func() {
			command, err := config.FindCommandConfig("not real")
			Expect(string(err.Error())).To(Equal("Command not found"))
			Expect(command.Name).To(Equal(""))
		})
	})

	Context("If we are looking for a subcommand from the config", func() {
		config := GetConfig()
		It("should return the subcommand config if found", func() {
			command, err := config.FindSubCommandConfig("metadata", "get-locations")
			Expect(command.Name).To(Equal("get-locations"))
			Expect(err).To(BeNil())
		})
		It("should return an error if subcommand not found", func() {
			command, err := config.FindSubCommandConfig("metadata", "not real")
			Expect(string(err.Error())).To(Equal("Subcommand not found"))
			Expect(command.Name).To(Equal(""))
		})
		It("should return an error if command not found", func() {
			command, err := config.FindSubCommandConfig("not real", "not real")
			Expect(string(err.Error())).To(Equal("Command not found"))
			Expect(command.Name).To(Equal(""))
		})
	})

	Context("If we need to generate an argument string", func() {
		Test = true
		config := GetConfig()
		It("should include all required args", func() {
			command, err := config.FindSubCommandConfig("test", "test1")
			Expect(err).To(BeNil())
			Expect(len(command.Args)).To(Equal(2))
			Expect(command.GetArgsUsage()).To(Equal("TEXT TEXT2 "))
		})
		It("should include optional args", func() {
			command, err := config.FindSubCommandConfig("test", "test2")
			Expect(err).To(BeNil())
			Expect(len(command.Args)).To(Equal(3))
			// probably not a good example, testing for order
			Expect(command.GetArgsUsage()).To(Equal("TEXT [OTHERTEXT] TEXT2 "))
		})
		It("should be empty if no args", func() {
			command, err := config.FindSubCommandConfig("test", "test3")
			Expect(err).To(BeNil())
			Expect(len(command.Args)).To(Equal(0))
			Expect(command.GetArgsUsage()).To(Equal(""))
		})
	})

	Context("If a flag is supplied check for new args and apply", func() {
		Test = true
		config := GetConfig()
		It("should include all required args", func() {
			command, err := config.FindSubCommandConfig("test", "test1")
			Expect(err).To(BeNil())
			Expect(len(command.Flags)).To(Equal(2))
			Expect(command.GetArgsUsage()).To(Equal("TEXT TEXT2 "))

			Expect(command.GetArgsUsageWithFlags("project")).To(Equal("PROJECT "))
		})

		It("should include all required args", func() {
			command, err := config.FindSubCommandConfig("test", "test1")
			Expect(err).To(BeNil())
			Expect(len(command.Args)).To(Equal(2))
			Expect(len(command.GetArgsForFlags("project"))).To(Equal(1))
			Expect(command.GetArgsForFlags("project").GetRequiredArgsCount()).To(Equal(1))
		})

		It("should include all flags with args", func() {
			command, err := config.FindSubCommandConfig("test", "test1")
			Expect(err).To(BeNil())
			Expect(len(command.Args)).To(Equal(2))
			Expect(len(command.GetFlagsWithArgs())).To(Equal(2))
			Expect(command.Args.GetRequiredArgsCount()).To(Equal(2))
		})
	})

	AfterEach(func() {
		Debug = false
		os.RemoveAll(ION_HOME)
	})
})
