package store

import (
  "encoding/json"
)

/**
 * UserCollection.
 */
type UserCollection struct {
  Users []User `json:"users"`
}

/**
 * Json()
 *
 * Return a string json representation of the user collection.
 */
func (userCollection *UserCollection) Json() (string, error) {
  out, err := json.Marshal(userCollection)
  if err != nil {
    return "", err
  }
  return string(out), nil
}

/**
 * GetEmptyUserCollection()
 *
 * Return an empty collection.
 */
func GetEmptyUserCollection() *UserCollection {
  return &UserCollection{[]User{}}
}
