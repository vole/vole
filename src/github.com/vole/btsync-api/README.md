# BTSync API

Golang client for the Bittorrent Sync API.

## Example

```go
package main

import (
  btsync "btsync-api"
  "fmt"
  "log"
)

func main() {
  api := btsync.New("login", "password", 8080, true)

  folders, err := api.GetFolders()
  if err != nil {
    log.Fatalf("Error! %s", err)
  }

  for _, folder := range *folders {
    fmt.Printf("Sync folder %s has %d files\n", folder.Dir, folder.Files)
  }

  speed, _ := api.GetSpeed()
  fmt.Printf("Speed: upload=%d, download=%d", speed.Upload, speed.Download)
}

```

## Documentation

http://godoc.org/github.com/vole/btsync-api
