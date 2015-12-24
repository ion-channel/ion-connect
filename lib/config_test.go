// config_test.go
//
// Copyright (C) 2015 Selection Pressure LLC
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package ionconnect

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "fmt"
  "os"
)

var _ = Describe("Config", func() {
  var (
  )

  BeforeEach(func() {
    Debug = false
    ION_HOME = fmt.Sprintf("%s/ionchannel-test", os.TempDir())
    CREDENTIALS_FILE = fmt.Sprintf("%s/credentials", ION_HOME)
  })

  Context("When the config.yaml file is loaded", func() {
    config := GetConfig()
    It("should contain an 'api' section", func() {
        Expect(len(config.Commands)).To(Equal(2))
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
        Expect(len(config.Commands[0].Subcommands)).To(Equal(3))
        Expect(config.Commands[0].Subcommands[0].Name).To(Equal("scan-git"))
        Expect(config.Commands[0].Subcommands[0].Post).To(BeTrue())
    })
    It("should have commands with subcommands and flags", func() {
        Expect(len(config.Commands[0].Subcommands)).To(Equal(3))
        Expect(config.Commands[0].Subcommands[0].Name).To(Equal("scan-git"))
        Expect(config.Commands[0].Subcommands[0].Flags[0].Name).To(Equal("project"))
    })
  })

  Context("When processing a url", func() {
    config := GetConfig()
    It("it should render template code", func() {
        params := GetParams{Id:"test"}
        url, err := config.ProcessUrlFromConfig("scanner", "get-scan", params)
        Expect(url).To(Equal("/scanner/getScan"))
        Expect(err).To(BeNil())
    })
    It("it should not fail if it's just a string", func() {
        params := GetParams{Id:"test"}
        url, err := config.ProcessUrlFromConfig("scanner", "scan-git", params)
        Expect(url).To(Equal("/scanner/scanGit"))
        Expect(err).To(BeNil())
    })
  })

  Context("When we need creds", func(){
    It("should create the ION_HOME directory", func() {
        Expect(LoadCredential()).To(Equal(""))
        Expect(PathExists(ION_HOME)).To(BeTrue())
    })
    It("should save credentials to new file", func() {
        Expect(LoadCredential()).To(Equal(""))
        Expect(PathExists(ION_HOME)).To(BeTrue())
        Expect(PathExists(CREDENTIALS_FILE)).To(BeFalse())
        Expect(func(){saveCredentials("notarealkey")}).ShouldNot(Panic())
        Expect(PathExists(CREDENTIALS_FILE)).To(BeTrue())
        Expect(LoadCredential()).To(Equal("notarealkey"))
    })
  })

  Context("If we are looking for a command from the config", func() {
    config := GetConfig()
    It("should return the command config if found", func() {
      command, err := config.FindCommandConfig("scanner")
      Expect(command.Name).To(Equal("scanner"))
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
      command, err := config.FindSubCommandConfig("scanner", "scan-git")
      Expect(command.Name).To(Equal("scan-git"))
      Expect(err).To(BeNil())
    })
    It("should return an error if subcommand not found", func() {
      command, err := config.FindSubCommandConfig("scanner", "not real")
      Expect(string(err.Error())).To(Equal("Subcommand not found"))
      Expect(command.Name).To(Equal(""))
    })
    It("should return an error if command not found", func() {
      command, err := config.FindSubCommandConfig("not real", "not real")
      Expect(string(err.Error())).To(Equal("Command not found"))
      Expect(command.Name).To(Equal(""))
    })
  })


  AfterEach(func() {
    Debug = false
    os.RemoveAll(ION_HOME)
  })
})
