package activitystreams

import (
  "encoding/json"
)

type MediaLink struct {
  Url      string `json:"string"`
  Duration int    `json:"duration"`
  Height   int    `json:"height"`
  Width    int    `json:"width"`
}

func (link *MediaLink) FromJson(j []byte) error {
  return json.Unmarshal(j, link)
}
