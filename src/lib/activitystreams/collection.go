package activitystreams

import (
  "encoding/json"
)

type Collection struct {
  TotalItems int      `json:"totalItems"`
  Items      []Object `json:"items"`
  Url        string   `json:"url"`
}

func (collection *Collection) FromJson(j []byte) error {
  return json.Unmarshal(j, collection)
}
