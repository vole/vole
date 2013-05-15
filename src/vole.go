package main

import (
  "net/http"
  "path/filepath"
  "fmt"
  "os"
  "strings"
  "io/ioutil"
)

func restHandler(w http.ResponseWriter, r *http.Request) {
  todos := make([]string, 0);

  filepath.Walk("./data", func(path string, info os.FileInfo, err error) error {
    if strings.HasPrefix(path, "data/todo-") {
      json, err := ioutil.ReadFile(path)
      if err == nil {
        todos = append(todos, string(json))
      }
    }
    return nil
  })
  out := `{ "posts": [`
  out += strings.Join(todos, ",")
  out += `] }`

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, out)
}

func main() {
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/" {
      http.ServeFile(w, r, "./web/index.html")
    } else {
      http.NotFound(w, r)
    }
  })
  http.HandleFunc("/rest/", restHandler)
  http.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("./js/"))))
  http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./css/"))))
  http.Handle("/img/", http.StripPrefix("/img", http.FileServer(http.Dir("./img/"))))
  http.ListenAndServe(":6789", nil)
}
