package config

import (
  "encoding/json"
  "errors"
  "fmt"
  "io/ioutil"
)

type Config struct {
  Install struct {
    Dir string `json:"dir"`
  } `json:"install"`

  UI struct {
    Reverse bool `json:"reverse"`
  } `json:"ui"`

  Server struct {
    Host string `json:"host"`
  } `json:"server"`
}

func Load() (*Config, error) {
  config := Config{}

  var file []byte
  file, err := ioutil.ReadFile("config.json")
  if err != nil {
    file, err = ioutil.ReadFile("config.sample.json")
    if err != nil {
      fmt.Println("Can't find config.json, reading default config.")
    }
  }

  err = json.Unmarshal(file, &config)
  if err != nil {
    return nil, errors.New("Unable to parse config.json.")
  }

  return &config, nil
}
