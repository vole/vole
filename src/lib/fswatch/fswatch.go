package main

import (
  "github.com/howeyc/fsnotify"
)

type fswatch struct {
  watchers map[string]fsnotify.Watcher
}

var fs = fswatch{
}

func (fs *fswatch) run() {
}
