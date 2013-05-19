package config

import (
  "errors"
  "io/ioutil"
  "launchpad.net/goyaml"
)

type Config struct {
  Install struct {
    Dir string `json:"dir"`
  } `json:"install"`

  UI struct {
    Reverse bool `json:"reverse"`
  } `json:"ui"`
}

func Load() (*Config, error) {
  config := Config{}

  file, err := ioutil.ReadFile("config.yaml")
  if err != nil {
    return nil, errors.New("Unable to open config.yaml. Make sure you copy config.sample.yaml to config.yaml.")
  }

  err = goyaml.Unmarshal(file, &config)
  if err != nil {
    return nil, errors.New("Unable to parse config.yaml.")
  }

  return &config, nil
}
