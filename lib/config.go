package ionconnect

import(
    "github.com/ion-channel/ion-connect/Godeps/_workspace/src/github.com/GeertJohan/go.rice"
    "github.com/ion-channel/ion-connect/Godeps/_workspace/src/gopkg.in/yaml.v2"
    "github.com/ion-channel/ion-connect/Godeps/_workspace/src/golang.org/x/crypto/ssh/terminal"
    "github.com/ion-channel/ion-connect/Godeps/_workspace/src/github.com/codegangsta/cli"
    "log"
    "fmt"
    "os"
    "errors"
    "text/template"
    "bytes"
)

type Command struct {
  Name string
  Usage string
  Post bool
  Url string
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


func (config Config) FindCommandConfig(commandName string) (Command, error) {
  for _, command := range config.Commands {
    if command.Name == commandName {
      return command, nil
    }
  }

  return Command{}, errors.New("Command not found")
}

func (config Config) ProcessUrlFromConfig(commandName string, subcommandName string, params interface{}) (string, error) {
  subCommandConfig, err := config.FindSubCommandConfig(commandName, subcommandName)
  if err != nil {
    return "", err
  }

  url := subCommandConfig.Url

  templ := template.Must(template.New("url").Parse(url))
  buf := bytes.Buffer{}
  err = templ.Execute(&buf, params)
  if err != nil {
    return "", err
  }

  return string(buf.Bytes()), nil
}

func (config Config) FindSubCommandConfig(commandName string, subcommandName string) (Command, error) {
  command, err := config.FindCommandConfig(commandName)
  if err != nil {
    return Command{}, err
  }

  for _, subcommand := range command.Subcommands {
    if subcommand.Name == subcommandName {
      return subcommand, nil
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
  secretKey, _ := terminal.ReadPassword(int(os.Stdin.Fd()))

  Debugf("All you keys are belong to us! (%s)", secretKey)

  if len(secretKey) != 0 {
    saveCredentials(string(secretKey))
  }
}

func LoadCredential() string {
  exists, _ := PathExists(ION_HOME)
  if exists {
    bytes, _ := ReadBytesFromFile(CREDENTIALS_FILE)
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
