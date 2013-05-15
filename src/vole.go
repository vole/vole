package main

import (
  "net/http"
  "fmt"
  "strings"
  "io/ioutil"
)

func getAllPosts() string {
  posts := make([]string, 0);

  users_info, _ := ioutil.ReadDir("./data/users")

  for _, user_info := range users_info {
    fmt.Println(string(user_info.Name()))
    posts_dir := "./data/users/" + string(user_info.Name()) + "/v1/posts"
    posts_info, _ := ioutil.ReadDir(posts_dir)

    for _, post_info := range posts_info {
      post_path := posts_dir + "/" + string(post_info.Name())
      fmt.Println(post_path)

      post_data, err := ioutil.ReadFile(post_path)
      if err == nil {
        posts = append(posts, string(post_data))
      }
    }
  }

  out := `{ "posts": [`
  out += strings.Join(posts, ",")
  out += `] }`

  fmt.Println(out)
  return out
}

func getMyUser() string {
  name := "";
  data, err := ioutil.ReadFile("data/my_user")
  if err == nil {
    name = strings.TrimSpace(string(data))
    fmt.Println("found my user: " + name)

    // Load the user's profile json.
    profile := ""
    profile_data, profile_err := ioutil.ReadFile("data/users/" + name + "/v1/user/" + name)

    if profile_err == nil {
      profile = strings.TrimSpace(string(profile_data))
    } else {
      fmt.Println(profile_err)
    }

    return `{ "users": [` + profile + `] }`
  }
  fmt.Println(err)
  return ""
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
  out := ""
  switch r.URL.Path {
    case "/api/posts":
      out = getAllPosts()
    case "/api/my_user":
      out = getMyUser()
  }

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
  http.HandleFunc("/api/", jsonHandler)
  http.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("./js/"))))
  http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./css/"))))
  http.Handle("/img/", http.StripPrefix("/img", http.FileServer(http.Dir("./img/"))))
  http.ListenAndServe(":6789", nil)
}
