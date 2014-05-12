package main

import (
	"github.com/vole/web"
	"lib/api"
	"lib/assets"
	"lib/config"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	logger := log.New(os.Stdout, "[Vole] ", log.Ldate|log.Ltime)
	web.SetLogger(log.New(os.Stdout, "[Web] ", log.Ldate|log.Ltime))

	logger.Printf("vole startup\n")

	web.Config.StaticDir = "./does-not-exist"

	web.Get("/status", api.Status)
	web.Get("/config", api.GetConfig)

	web.Get("/api/drafts", api.GetDrafts)
	web.Get("/api/drafts/(.*)", api.GetDraft)
	web.Post("/api/drafts", api.CreateDraft)
	web.Put("/api/drafts/(.*)", api.SaveDraft)
	web.Delete("/api/drafts/(.*)", api.DeleteDraft)

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
