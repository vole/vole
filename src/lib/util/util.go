package util

import (
  "errors"
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

func Load() (*Config, error) {
  config := &Config{}

  var file []byte
  file, err := ioutil.ReadFile("config.json")
  if err != nil {
    file, err = ioutil.ReadFile("config.default.json")
    if err != nil {
      fmt.Println("Can't find config.json, reading default config.")
    }
  }

  err = json.Unmarshal(file, config)
  if err != nil {
    return nil, errors.New("Unable to parse config.json.")
  }

  return config, nil
}
