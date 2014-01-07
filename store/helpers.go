package store

import (
  "io/ioutil"
  "os"
  "path"
)

/**
 * Helper functions.
 */
func ReadFile(args ...string) ([]byte, error) {
  return ioutil.ReadFile(path.Join(args...))
}

func ReadDir(args ...string) ([]os.FileInfo, error) {
  return ioutil.ReadDir(path.Join(args...))
}

func Create(args ...string) (*os.File, error) {
  return os.Create(path.Join(args...))
}

func Write(path string, data []byte) error {
  file, err := os.Create(path)
  if err != nil {
    return err
  }

  _, err = file.Write(data)
  if err != nil {
    return err
  }

  return file.Close()
}

func Delete(args ...string) error {
  return os.Remove(path.Join(args...))
}
