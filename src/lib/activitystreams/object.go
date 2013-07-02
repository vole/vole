package activitystreams

import (
  "encoding/json"
  "time"
)

type Object struct {
  Id                   string    `json:"id"`
  Url                  string    `json:"url"`
  Attachments          *[]Object `json:"attachments"`
  Author               *Object   `json:"author"`
  Content              string    `json:"content"`
  DisplayName          string    `json:"displayName"`
  DownstreamDuplicates []string  `json:"downstreamDuplicates"`
  Image                MediaLink `json:"image"`
  ObjectTyoe           string    `json:"objectType"`
  Summary              string    `json:"summary"`
  Updated              time.Time `json:"updated"`
  Published            time.Time `json:"published"`
  UpstreamDuplicates   []string  `json:"upstreamDuplicates"`
}

func (object *Object) FromJson(j []byte) error {
  return json.Unmarshal(j, object)
}
