package main

import (
  "fmt"
  "flag"
  "strings"
  "io/ioutil"
  "github.com/hoisie/web"
)

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

  //fmt.Println(out)

//   out = `{ "posts": [{
//   "id": 1,
//   "title": "post number 1",
//   "user": "billy"
// }
// ,{
//   "id": 2,
//   "title": "post number 2",
//   "user": "billy"
// }
// ,{
//   "id": 3,
//   "title": "post number 3",
//   "user": "billy"
// }]}`

  return out
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
  web.Get("/api/my_user", getMyUser)
  web.Run("0.0.0.0:" + *port)
}
