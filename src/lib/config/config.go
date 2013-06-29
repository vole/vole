package config

import (
  "encoding/json"
  "errors"
  "io/ioutil"
)

type Config struct {
  Install struct {
    Dir string `json:"dir"`
  } `json:"install"`

  UI struct {
    Reverse      bool `json:"reverse"`
    PollInterval int  `json:"pollInterval"`
  } `json:"ui"`

  Server struct {
    Listen string `json:"listen"`
    Debug  bool   `json:"debug"`
  } `json:"server"`
}

func Load() (*Config, error) {
  // Create configuration object and set default values.
  // Make sure any changes here are reflected in config.sample.json,
  // so that users who copy the file to modify defaults don't have old
  // values.
  config := Config{}
  config.Install.Dir = "~/Vole"
  config.UI.Reverse = false
  config.UI.PollInterval = 5000
  config.Server.Listen = "127.0.0.1:6789"
  config.Server.Debug = false

  // Now read config.json for any overrides of the defaults.
  var file []byte
  file, err := ioutil.ReadFile("config.json")
  if err == nil {
    err = json.Unmarshal(file, &config)
    if err != nil {
      return nil, errors.New("Unable to parse config.json.")
    }
  }

  return &config, nil
}
