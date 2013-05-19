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
    Reverse bool `json:"reverse"`
  } `json:"ui"`
}

func Load() (*Config, error) {
  config := Config{}

  file, err := ioutil.ReadFile("config.json")
  if err != nil {
    return nil, errors.New("Unable to open config.json. Make sure you copy config.sample.json to config.json.")
  }

  err = json.Unmarshal(file, &config)
  if err != nil {
    return nil, errors.New("Unable to parse config.json.")
  }

  return &config, nil
}
