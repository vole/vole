package store

import (
  "encoding/json"
)

/**
 * PostContainer.
 */
type PostContainer struct {
  Post Post `json:"post"`
}

/**
 * Json()
 *
 * Return a string json representation of the post container.
 */
func (postContainer *PostContainer) Json() (string, error) {
  out, err := json.Marshal(postContainer)
  if err != nil {
    return "", err
  }
  return string(out), nil
}
