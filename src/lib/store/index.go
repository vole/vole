package store

import (
  "encoding/json"
  "fmt"
  "fsnotify"
  "os"
  "path"
  "path/filepath"
  "strings"
)

type Index struct {
  Posts []IndexItem
}

type IndexItem struct {
  Path     string
  Created  int64
  Modified int64
}

func FileIsPost(file string) bool {
  ext := path.Ext(file)
  if ext != ".json" {
    return false
  }

  return !strings.HasSuffix(file, "user.json")
}

var instance *Index = nil

func BuildIndex(root string) (*Index, error) {
  fmt.Println("building index")

  watcher, err := inotify.NewWatcher()

  items := make([]IndexItem, 0)

  visit := func(file string, f os.FileInfo, err error) error {
    if !FileIsPost(file) {
      return nil
    }

    data, err := ReadFile(file)
    if err != nil {
      return nil
    }

    var post Post
    if err := json.Unmarshal(data, &post); err != nil {
      return nil
    }

    stat, err := os.Stat(file)
    if err != nil {
      return nil
    }

    item := &IndexItem{Path: file, Created: post.Created, Modified: stat.ModTime().UnixNano()}

    items = append(items, *item)
    return nil
  }

  err := filepath.Walk(root, visit)
  if err != nil {
    return nil, err
  }

  instance = &Index{Posts: items}

  fmt.Println("done.")

  return instance, nil
}
