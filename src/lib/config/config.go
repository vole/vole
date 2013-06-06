package config

import (
  "encoding/json"
  "errors"
  "fmt"
  "io/ioutil"
  osuser "os/user"
  "strings"
)

type Config struct {
  Install struct {
    Dir string `json:"dir"`
  } `json:"install"`

  UI struct {
    Reverse bool   `json:"reverse"`
    Theme   string `json:"theme"`
  } `json:"ui"`
}

func Load() *Config {
  config := &Config{}

  var file []byte
  file, err := ioutil.ReadFile("config.json")
  if err != nil {
    fmt.Println("Can't find config.json, reading default config.")
    file, err = ioutil.ReadFile("config.default.json")

    if err != nil {
      panic(errors.New("Unable to find config.json"))
    }
  }

  err = json.Unmarshal(file, config)
  if err != nil {
    panic(errors.New("Unable to parse config.json."))
  }

  return config
}

func (config *Config) InstallDir() string {
  conf := Load()

  user, err := osuser.Current()
  if err != nil {
    return conf.Install.Dir
  }

  dir := strings.Replace(conf.Install.Dir, "{HOME}", user.HomeDir, 1)

  if strings.HasPrefix(dir, "~") {
    dir = strings.Replace(conf.Install.Dir, "~", user.HomeDir, 1)
  }

  return dir
}
