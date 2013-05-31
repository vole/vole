package main

import (
  //"fmt"
  "encoding/json"
  "flag"
  "github.com/vole/web"
  "io/ioutil"
  "lib/config"
  "lib/store"
  osuser "os/user"
  "path"
)

var port = flag.String("port", "6789", "Port on which to run the web server.")

var DIR = func() string {
  dir := "."
  user, err := osuser.Current()
  if err == nil {
    dir = user.HomeDir
  }
  return path.Join(dir, "Vole")
}()

var userStore = &store.UserStore{
  Path:    DIR,
  Version: "v1",
}

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

    var allPosts *store.PostCollection
    allPosts, err := userStore.GetPosts()
    if err != nil || len(allPosts.Posts) < 1 {
      // Return a welcome post.
      post := &store.Post{}
      post.InitNew("Welcome to Vole. To start, create a new profile", "none", "none", "Welcome", "", false)
      post.Id = "none"
      allPosts = post.Collection()
    }

    postsJson, err := allPosts.Json()
    if err != nil {
      ctx.Abort(500, "Error getting posts as json.")
    }

    return postsJson
  })

  web.Get("/api/users", func(ctx *web.Context) string {
    ctx.ContentType("json")

    _, isMyUserFilter := ctx.Params["is_my_user"]

    var users *store.UserCollection

    if isMyUserFilter {
      myUser, _ := userStore.GetMyUser()
      if myUser != nil {
        users = myUser.Collection()
      } else {
        users = store.GetEmptyUserCollection()
      }
    } else {
      users, err = userStore.GetUsers()
      if err != nil {
        ctx.Abort(500, "Error loading all users.")
      }
    }

    usersJson, err := users.Json()
    if err != nil {
      ctx.Abort(500, "Error getting users as json.")
    }

    return usersJson
  })

  web.Post("/api/users", func(ctx *web.Context) string {
    body, err := ioutil.ReadAll(ctx.Request.Body)
    if err != nil {
      ctx.Abort(500, "Error reading request body.")
    }
    user, err := userStore.NewUserFromContainerJson(body)
    if err != nil {
      ctx.Abort(500, "Invalid JSON")
    }
    if err := user.Save(); err != nil {
      ctx.Abort(500, "Error saving user")
    }
    if err := userStore.SetMyUser(user); err != nil {
      ctx.Abort(500, "Error setting my user")
    }

    container := user.Container()
    userJson, err := container.Json()
    if err != nil {
      ctx.Abort(500, "Could not create container")
    }
    return userJson
  })

  web.Post("/api/posts", func(ctx *web.Context) string {
    body, err := ioutil.ReadAll(ctx.Request.Body)
    if err != nil {
      ctx.Abort(500, "Error reading request body.")
    }

    user, err := userStore.GetMyUser()
    if err != nil {
      ctx.Abort(500, "Error reading my user when posting.")
    }
    post, err := user.NewPostFromContainerJson(body)
    if err != nil {
      ctx.Abort(500, "Invalid JSON")
    }
    if err := post.Save(); err != nil {
      ctx.Abort(500, "Error saving post")
    }
    container := post.Container()
    postJson, err := container.Json()
    if err != nil {
      ctx.Abort(500, "Could not create container")
    }
    return postJson
  })

  web.Delete("/api/posts/(.*)", func(ctx *web.Context, id string) string {
    user, err := userStore.GetMyUser()
    if err != nil {
      ctx.Abort(500, "Error loading user.")
    }

    posts, err := user.GetPosts()
    if err != nil {
      ctx.Abort(500, "Error loading posts.")
    }

    for _, post := range posts.Posts {
      if post.Id == id {
        err := post.Delete()
        if err != nil {
          ctx.Abort(500, "Error deleting post.")
        } else {
          return "OK"
        }
      }
    }

    return "OK"
  })

  web.Run("0.0.0.0:" + *port)
}
