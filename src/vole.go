package main

import (
  "os"
  "fmt"
  "flag"
  "strings"
  "io/ioutil"
  "time"
  "github.com/vole/web"
  "github.com/nu7hatch/gouuid"
  "encoding/json"
)

/**
 *
 * import "vole/db"
 *
 * db.findPosts(group)
 * db.findUsers(group)
 *
 * db.findUser(group, KEY)
 * db.findPost(group, GUID)
 *
 * db.savePost(group, post)
 *
 */

func getAllPosts(ctx *web.Context) string {
  ctx.SetHeader("Content-Type", "application/json", true)

  posts := make([]string, 0)

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

  return out
}

func savePost(ctx *web.Context) string {
  data, err := ioutil.ReadFile("data/my_user")
  if err != nil {
    ctx.Abort(500, "Couldn't determine current user.")
  }
  //var jsonBlob = []byte(`{"post":{"title":"hello","user":"mark"}}`)

  body, err := ioutil.ReadAll(ctx.Request.Body);

  var post interface{}
  err = json.Unmarshal(body, &post)
  if err != nil {
    fmt.Println("error:", err)
  }

  user := strings.TrimSpace(string(data))

  ts := time.Now().UnixNano()
  uuid, _ := uuid.NewV4()
  filename := fmt.Sprintf("%d-post-%s", ts, uuid)

  file, err := os.Create("data/users/" + user + "/v1/posts/" + filename)
  if err != nil {
    ctx.Abort(500, "Unable to create file.")
  }

  m := post.(map[string]interface{})

  m["post"] = "dkjfbsdjf"

  ohnoes = m.(map[string]interface{})

  lol, _ := json.Marshal(post)
  fmt.Printf("%+v", string(lol))

  file.Write(lol)
  return string(lol)
}

func getMyUser(ctx *web.Context) string {
  ctx.SetHeader("Content-Type", "application/json", true)

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

func main() {
  var port = flag.String("port", "6789", "Port on which to run the web server.")
  flag.Parse()

  web.Get("/api/posts", getAllPosts)
  web.Get("/api/users", getMyUser)

  web.Post("/api/posts", savePost)

  web.Run("0.0.0.0:" + *port)
}
