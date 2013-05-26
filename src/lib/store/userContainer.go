package store

import (
  "encoding/json"
)

/**
 * UserContainer.
 */
type UserContainer struct {
  User User `json:"user"`
}

/**
 * Json()
 *
 * Return a string json representation of the user container.
 */
func (userContainer *UserContainer) Json() (string, error) {
  out, err := json.Marshal(userContainer)
  if err != nil {
    return "", err
  }
  return string(out), nil
}
