package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/petar/GoLLRB/llrb"
	"github.com/vole/web"
	"io/ioutil"
)

var port = flag.String("port", "6789", "Port on which to run the web server.")

func main() {
	flag.Parse()

	// tree := GlobalPostTree()
	// tree.AscendGreaterOrEqual(tree.Min(), func(item llrb.Item) bool {
	// 	i, ok := item.(llrb.Int)
	// 	if !ok {
	// 		return false
	// 	}
	// 	fmt.Println(int(i))
	// 	return true
	// })

	web.Get("/api/posts", func(ctx *web.Context) string {
		ctx.ContentType("json")

		posts, err := getPosts()
		if err != nil {
			ctx.Abort(500, "Error loading posts.")
		}

		postsJson, err := json.Marshal(*posts)
		if err != nil {
			ctx.Abort(500, "Error marshalling posts.")
		}

		return string(postsJson)
	})

	web.Get("/api/users", func(ctx *web.Context) string {
		ctx.ContentType("json")

		user, err := CurrentUser()
		if err != nil {
			ctx.Abort(500, "Error loading user.")
		}

		collection := NewUserCollection([]User{*user})

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

		container := PostContainerFromJson(body)

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
