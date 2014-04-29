package config

import (
	"encoding/json"
	"io/ioutil"
	osuser "os/user"
	"path"
	"reflect"
)

type Config struct {
	// Install configuration.
	Install_Dir string `json:"install_dir"`

	// UI configuration.
	UI_Layout       string `json:"ui_layout"`
	UI_Theme        string `json:"ui_theme"`
	UI_Reverse      bool   `json:"ui_reverse"`
	UI_PollInterval int    `json:"ui_pollInterval"`
	UI_PageSize     int    `json:"ui_pageSize"`

	// Server configuration.
	Server_Listen string `json:"server_listen"`
	Server_Debug  bool   `json:"server_debug"`
}

var config = func() *Config {
	// Create configuration object and set default values.
	// Make sure any changes here are reflected in config.sample.json,
	// so that users who copy the file to modify defaults don't have old
	// values.
	config := Config{}
	config.Install_Dir = "~/Vole"
	config.UI_Layout = "default"
	config.UI_Theme = "default"
	config.UI_Reverse = false
	config.UI_PollInterval = 5000
	config.Server_Listen = "127.0.0.1:6789"
	config.Server_Debug = false
	config.UI_PageSize = 50

	// Now read config.json for any overrides of the defaults.
	file, err := ioutil.ReadFile("config.json")
	if err == nil {
		err = json.Unmarshal(file, &config)
		if err != nil {
			panic(err)
		}
	}

	return &config
}()

func Json() (string, error) {
	bytes, err := json.Marshal(config)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func Version() string {
	return "v1"
}

func StorageDir() string {
	dir := "."
	user, err := osuser.Current()
	if err == nil {
		dir = user.HomeDir
	}
	return path.Join(dir, "Vole")
}

func ReadString(field string) string {
	r := reflect.ValueOf(config)
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.String())
}

func ReadInt(field string) int {
	r := reflect.ValueOf(config)
	f := reflect.Indirect(r).FieldByName(field)
	return int(f.Int())
}

func ReadBool(field string) bool {
	r := reflect.ValueOf(config)
	f := reflect.Indirect(r).FieldByName(field)
	return bool(f.Bool())
}
