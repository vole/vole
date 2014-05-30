package main

import (
	"flag"
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

var configFile = flag.String("config", "config.json", "Path to the Vole config file.")

func main() {
	flag.Parse()
	config.Load(*configFile)

	logger := log.New(os.Stdout, "[Vole] ", log.Ldate|log.Ltime)
	web.SetLogger(log.New(os.Stdout, "[Web]  ", log.Ldate|log.Ltime))

	logger.Printf("Server is starting\n")

	web.Config.StaticDir = "./does-not-exist"

	web.Get("/status", api.GetStatus)
	web.Get("/config", api.GetConfig)

	web.Get("/me", api.GetVoleUser)
	web.Post("/me", api.CreateVoleUser)
	web.Put("/me", api.UpdateVoleUser)
	web.Delete("/me", api.DeleteVoleUser)

	web.Get("/api/drafts", api.GetDrafts)
	web.Get("/api/drafts/(.*)", api.GetDraft)
	web.Post("/api/drafts", api.CreateDraft)
	web.Put("/api/drafts/(.*)", api.SaveDraft)
	web.Delete("/api/drafts/(.*)", api.DeleteDraft)

	web.Get("/api/posts", api.GetPosts)
	web.Get("/api/posts/(.*)", api.GetPost)
	web.Post("/api/posts", api.CreatePost)
	web.Delete("/api/posts/(.*)", api.DeletePost)

	web.Get("/api/users", api.GetUsers)
	web.Get("/api/users/(.*)", api.GetUser)
	web.Put("/api/users/(.*)", api.SaveUser)
	web.Delete("/api/users/(.*)", api.DeleteUser)

	// Serve static files using the asset manifest.
	web.Get("/.*", func(ctx *web.Context) []byte {
		filePath := strings.TrimLeft(ctx.Request.URL.Path, "/")

		file, err := assets.Asset(filePath)
		if err != nil {
			ctx.SetHeader("Content-Security-Policy", "script-src 'self' 'unsafe-eval'", true)
			ctx.SetHeader("Content-Type", "text/html", true)
			file, _ = assets.Asset("index.html")
		} else {
			mimeType := mime.TypeByExtension(filepath.Ext(filePath))
			ctx.SetHeader("Content-Type", mimeType, true)
		}

		return file
	})

	web.Run(config.ReadString("Server_Listen"))
}
