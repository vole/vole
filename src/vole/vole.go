package main

import (
	"fmt"
	"github.com/vole/web"
	"lib/api"
	"lib/assets"
	"lib/config"
	"mime"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("vole startup")

	web.Config.StaticDir = "./does-not-exist"

	web.Get("/status", api.Status)
	web.Get("/config", api.GetConfig)
	web.Get("/api/drafts", api.GetDrafts)
	web.Get("/api/drafts/(.*)", api.GetDraft)
	web.Get("/api/posts", api.GetPosts)
	web.Get("/api/posts/(.*)", api.GetPost)
	web.Post("/api/posts", api.SavePost)
	web.Delete("/api/posts/(.*)", api.DeletePost)
	web.Get("/api/users", api.GetUsers)
	web.Post("/api/users", api.SaveUser)
	web.Post("/api/friend", api.SaveFriend)

	web.Get("/.*", func(ctx *web.Context) []byte {
		filePath := strings.TrimLeft(ctx.Request.URL.Path, "/")

		file, err := assets.Asset(filePath)
		if err != nil {
			ctx.SetHeader("Content-Security-Policy", "script-src 'self' 'unsafe-eval'", true)
			ctx.SetHeader("Content-Type", "text/html", true)
			file, _ = assets.Asset("index.html")
		} else {
			fileExt := filepath.Ext(filePath)
			mimeType := mime.TypeByExtension(fileExt)
			ctx.SetHeader("Content-Type", mimeType, true)
		}

		return file
	})

	web.Run(config.ReadString("Server_Listen"))
}
