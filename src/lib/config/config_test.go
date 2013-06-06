package config

import (
  "testing"
)

func TestInstallDir(t *testing.T) {
  conf := Load()
  conf.InstallDir()
}
