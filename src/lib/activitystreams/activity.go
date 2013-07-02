package activitystreams

import (
  "encoding/json"
  "time"
)

type Activity struct {
  Id        string    `json:"id"`
  Title     string    `json:"title"`
  Content   string    `json:"content"`
  Url       string    `json:"url"`
  Actor     Object    `json:"actor"`
  Verb      string    `json:"verb"`
  Object    Object    `json:"object"`
  Target    Object    `json:"target"`
  Generator Object    `json:"generator"`
  Provider  Object    `json:"provider"`
  Icon      MediaLink `json:"icon"`
  Updated   time.Time `json:"updated"`
  Published time.Time `json:"published"`
}

func (activity *Activity) FromJson(j []byte) error {
  return json.Unmarshal(j, activity)
}
