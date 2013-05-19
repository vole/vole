package main

import (
  "encoding/json"
  "flag"
  "github.com/vole/web"
  "io/ioutil"
  "lib/config"
  "lib/db"
  "sort"
)

var port = flag.String("port", "6789", "Port on which to run the web server.")

func main() {
  flag.Parse()

  config, err := config.Load()
  if err != nil {
    panic(err)
  }

  web.Get("/api/config", func(ctx *web.Context) string {
    ctx.ContentType("json")

    configJson, err := json.Marshal(config)
    if err != nil {
      ctx.Abort(500, "Error marshalling config.")
    }

    return string(configJson)
  })

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

    var collection *db.UserCollection

    _, isMyUserFilter := ctx.Params["is_my_user"]

    if isMyUserFilter {
      currentUser, _ := db.CurrentUser()
      if currentUser != nil {
        collection = db.NewUserCollection([]db.User{*currentUser})
      } else {
        collection = db.NewUserCollection([]db.User{})
      }
    } else {
      users, err := db.GetUsers()
      if err != nil {
        ctx.Abort(500, "Error loading users.")
      }
      collection = users
    }

    usersJson, err := json.Marshal(collection)
    if err != nil {
      ctx.Abort(500, "Error marshalling users.")
    }

    return string(usersJson)
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
