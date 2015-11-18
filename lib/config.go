package ionconnect

import(
    "github.com/ion-channel/ion-connect/Godeps/_workspace/src/github.com/GeertJohan/go.rice"
    "github.com/ion-channel/ion-connect/Godeps/_workspace/src/gopkg.in/yaml.v2"
    "github.com/ion-channel/ion-connect/Godeps/_workspace/src/golang.org/x/crypto/ssh/terminal"
    "github.com/ion-channel/ion-connect/Godeps/_workspace/src/github.com/codegangsta/cli"
    "log"
    "fmt"
    "syscall"
    "errors"
)

type Command struct {
  Name string
  Usage string
  Write bool
  Flags []Flag
  Subcommands []Command
}

type Flag struct {
  Name string
  Value string
  Usage string
}

type Config struct {
  Version string
  Endpoint string
  Token string
  Commands []Command
}

func GetConfig() Config {
  configBox, err := rice.FindBox("../config")
  if err != nil {
    log.Fatal(err)
  }

  // get file contents as string
  configString, err := configBox.String("config.yaml")
  if err != nil {
    log.Fatal(err)
  }

  config := Config{}

  err = yaml.Unmarshal([]byte(configString), &config)
  if err != nil {
    log.Fatalf("error: %v", err)
  }
  Debugf("Config map:\n%v\n\n", config)

  return config
}


func (config Config) findCommandConfig(commandName string) (Command, error) {
  for index := range config.Commands {
    if config.Commands[index].Name == commandName {
      return config.Commands[index], nil
    }
  }

  return Command{}, errors.New("Command not found")
}

func (config Config) findSubCommandConfig(commandName string, subcommandName string) (Command, error) {
  command, err := config.findCommandConfig(commandName)
  if err != nil {
    return Command{}, err
  }

  for index := range command.Subcommands {
    if command.Subcommands[index].Name == subcommandName {
      return command.Subcommands[index], nil
    }
  }

  return Command{}, errors.New("Subcommand not found")
}


func HandleConfigure(context* cli.Context) {
  currentSecretKey := LoadCredential()
  truncatedSecretKey := currentSecretKey
  if len(currentSecretKey) > 4 {
      truncatedSecretKey = currentSecretKey[len(currentSecretKey)-4:len(currentSecretKey)]
  }

  fmt.Printf("Ion Channel Api Key [%s]: ", truncatedSecretKey)
  secretKey, _ := terminal.ReadPassword(syscall.Stdin)

  Debugf("All you keys are belong to us! (%s)", secretKey)

  if len(secretKey) != 0 {
    saveCredentials(string(secretKey))
  }
}

func LoadCredential() string {
  exists, _ := PathExists(ION_HOME)
  if exists {
    bytes, _ := ReadBytesFromFile(CREDENTIALS_FILE)
    Debugf("Reading credentials file (%s)", bytes)
    credentials := make(map[string]string)
    yaml.Unmarshal([]byte(bytes), &credentials)
    return credentials[CREDENTIALS_KEY_FIELD]
  } else {
    MkdirAll(ION_HOME, 0775)
    return ""
  }
}

func saveCredentials(secretKey string) {
  credentials := make(map[string]string)
  credentials[CREDENTIALS_KEY_FIELD] = secretKey
  yamlCredentials, _ := yaml.Marshal(&credentials)
  WriteLinesToFile(CREDENTIALS_FILE, []string{string(yamlCredentials)}, 0600)
}
