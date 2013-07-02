package activitystreams

import (
  "io/ioutil"
  "reflect"
  "testing"
  "time"
)

func TestActivityFromJson(t *testing.T) {
  json, err := ioutil.ReadFile("test/activity.json")
  if err != nil {
    t.Errorf("Error reading file for tests: %+v", err)
  }

  activity := Activity{}

  err = activity.FromJson(json)
  if err != nil {
    t.Errorf("Error unmarshalling activity: %+v", err)
  }

  if activity.Actor.DisplayName != "Martin Smith" {
    t.Error("Expected Actor.DiplayName to equal 'Martin Smith'")
  }

  if activity.Target.DisplayName != "Martin's Blog" {
    t.Error("Expected Target.DisplayName to equal 'Martin's Blog'")
  }

  published := time.Date(2011, time.February, 10, 15, 4, 55, 0, time.UTC)
  if activity.Published != published {
    t.Errorf("Expected Published to equal %+v, was: %+v", published, activity.Published)
  }

  actorType := reflect.TypeOf(activity.Actor)
  if actorType.Name() != "Object" {
    t.Errorf("Expected Actor to be an Object, was: %+v", actorType)
  }

}
