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
        Expect(len(config.Commands)).To(Equal(1))
    })
    It("should contain the api version", func() {
        Expect(config.Version).To(Equal("v1"))
    })
    It("should contain the api endpoint", func() {
        Expect(config.Endpoint).To(Equal("https://api.ionchannel.io/"))
    })
    It("should contain the api token header name", func() {
        Expect(config.Token).To(Equal("access-token"))
    })
    It("should have commands with subcommands", func() {
        Expect(len(config.Commands[0].Subcommands)).To(Equal(1))
        Expect(config.Commands[0].Subcommands[0].Name).To(Equal("scan"))
        Expect(config.Commands[0].Subcommands[0].Write).To(BeTrue())
    })
    It("should have commands with subcommands and flags", func() {
        Expect(len(config.Commands[0].Subcommands)).To(Equal(1))
        Expect(config.Commands[0].Subcommands[0].Name).To(Equal("scan"))
        Expect(config.Commands[0].Subcommands[0].Flags[0].Name).To(Equal("name"))
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
      command, err := config.findCommandConfig("scanner")
      Expect(command.Name).To(Equal("scanner"))
      Expect(err).To(BeNil())
    })
    It("should return an error if not found", func() {
      command, err := config.findCommandConfig("not real")
      Expect(string(err.Error())).To(Equal("Command not found"))
      Expect(command.Name).To(Equal(""))
    })
  })

  Context("If we are looking for a subcommand from the config", func() {
    config := GetConfig()
    It("should return the subcommand config if found", func() {
      command, err := config.findSubCommandConfig("scanner", "scan")
      Expect(command.Name).To(Equal("scan"))
      Expect(err).To(BeNil())
    })
    It("should return an error if subcommand not found", func() {
      command, err := config.findSubCommandConfig("scanner", "not real")
      Expect(string(err.Error())).To(Equal("Subcommand not found"))
      Expect(command.Name).To(Equal(""))
    })
    It("should return an error if command not found", func() {
      command, err := config.findSubCommandConfig("not real", "not real")
      Expect(string(err.Error())).To(Equal("Command not found"))
      Expect(command.Name).To(Equal(""))
    })
  })


  AfterEach(func() {
    Debug = false
    os.RemoveAll(ION_HOME)
  })
})
