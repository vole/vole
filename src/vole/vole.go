package main

import (
  "encoding/json"
  "flag"
  "github.com/vole/web"
  "io/ioutil"
  "lib/db"
  "sort"
)

var port = flag.String("port", "6789", "Port on which to run the web server.")

func main() {
  flag.Parse()

  web.Get("/api/posts", func(ctx *web.Context) string {
    ctx.ContentType("json")

    posts, err := db.GetPosts()
    if err != nil {
      ctx.Abort(500, "Error loading posts.")
    }

    sort.Sort(posts)

    postsJson, err := json.Marshal(*posts)
    if err != nil {
      ctx.Abort(500, "Error marshalling posts.")
    }

    return string(postsJson)
  })

  web.Get("/api/users", func(ctx *web.Context) string {
    ctx.ContentType("json")

    user, err := db.CurrentUser()
    if err != nil {
      ctx.Abort(500, "Error loading user.")
    }

    collection := db.NewUserCollection([]db.User{*user})

    userJson, err := json.Marshal(collection)
    if err != nil {
      ctx.Abort(500, "Error marshalling user.")
    }

    return string(userJson)
  })

  web.Post("/api/posts", func(ctx *web.Context) string {
    body, err := ioutil.ReadAll(ctx.Request.Body)
    if err != nil {
      ctx.Abort(500, "Error reading request body.")
    }

    container := db.PostContainerFromJson(body)

    err = (*container).Post.Save()
    if err != nil {
      ctx.Abort(500, "Error saving post.")
    }

    containerJson, err := json.Marshal(*container)
    if err != nil {
      ctx.Abort(500, "Error marshalling post.")
    }

    return string(containerJson)
  })

  web.Run("0.0.0.0:" + *port)
}
