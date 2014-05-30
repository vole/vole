package api

import (
	"encoding/json"
	"fmt"
	btsync "github.com/vole/btsync-api"
	"github.com/vole/gouuid"
	"github.com/vole/web"
	"io/ioutil"
	"lib/config"
	"lib/store"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

type Error struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

var (
	logger = log.New(os.Stdout, "[Vole] ", log.Ldate|log.Ltime)
)

// Set the correct HTTP headers for a JSON response.
func setJsonHeaders(ctx *web.Context) {
	ctx.ContentType("json")
	ctx.SetHeader("Cache-Control", "no-cache, no-store", true)
}

// Create a JSON error response.
func createJsonError(ctx *web.Context, status int, message string) string {
	ctx.ResponseWriter.WriteHeader(status)
	var err = Error{true, message}
	bytes, _ := json.Marshal(err)
	return string(bytes)
}

// Get the current status of the backend. Intended to report
// connectivity to various services, like BTSync, etc.
func GetStatus(ctx *web.Context) string {
	setJsonHeaders(ctx)

	api := btsync.New(
		config.ReadString("BTSync_User"),
		config.ReadString("BTSync_Pass"),
		config.ReadInt("BTSync_Port"),
		true)

	_, err := api.GetVersion()
	if err != nil {
		return "{ \"btsync\": false }"
	}

	return "{ \"btsync\": true }"
}

// Fetch the app config.
// TODO(aaron): Remove sensitive data from the config
// before sending it to the UI.
func GetConfig(ctx *web.Context) string {
	setJsonHeaders(ctx)

	configJson, err := config.Json()
	if err != nil {
		logger.Printf("%s", err)
		return createJsonError(ctx, 500, "Error loading config.")
	}

	return configJson
}

// Create the Vole user.
func CreateVoleUser(ctx *web.Context) string {
	setJsonHeaders(ctx)

	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		logger.Printf("%s", err)
		return createJsonError(ctx, 500, "Error reading request body.")
	}

	var user = &store.User{}
	if err := json.Unmarshal(body, user); err != nil {
		logger.Printf("%s", err)
		return createJsonError(ctx, 500, "Error unmarshalling user.")
	}

	if err := store.Load().CreateVoleUser(user); err != nil {
		logger.Printf("%s", err)
		return createJsonError(ctx, 500, "Error saving user.")
	}

	userJson, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		logger.Printf("%s", err)
		return createJsonError(ctx, 500, "Error marshalling user.")
	}

	return string(userJson)
}

// Get the current Vole user.
func GetVoleUser(ctx *web.Context) string {
	setJsonHeaders(ctx)

	user, err := store.Load().GetVoleUser()
	if err != nil {
		logger.Printf("%s", err)
		return createJsonError(ctx, 404, "User not found.")
	}

	userJson, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		logger.Printf("%s", err)
		return createJsonError(ctx, 500, "Error marshalling user.")
	}

	return string(userJson)
}

// Update the current Vole user.
// TODO(aaron): UpdateVoleUser()
func UpdateVoleUser(ctx *web.Context) string {
	return ""
}

// Delete the current Vole user.
// TODO(aaron): DeleteVoleUser()
func DeleteVoleUser(ctx *web.Context) string {
	return ""
}

// Get all posts. Available parameters:
// before - Get posts before the specified time.
// limit - Maximum number of posts to return.
// user - Only get the specific user's posts.
func GetPosts(ctx *web.Context) string {
	setJsonHeaders(ctx)

	before, _ := ctx.Params["before"]
	userId, _ := ctx.Params["user"]

	posts, err := store.Load().GetPosts(userId, before, config.ReadInt("UI_PageSize"))
	if err != nil {
		return createJsonError(ctx, 500, fmt.Sprintf("%s", err))
	}

	postsJson, err := json.MarshalIndent(posts, "", "  ")
	if err != nil {
		return createJsonError(ctx, 500, "Error getting posts as json.")
	}

	return string(postsJson)
}

// Get a list of all users.
func GetUsers(ctx *web.Context) string {
	setJsonHeaders(ctx)

	query, hasQuery := ctx.Params["query"]

	users, err := store.Load().GetUsers()
	if err != nil {
		logger.Printf("%s", err)
		return createJsonError(ctx, 500, "Error getting users.")
	}

	if hasQuery {
		users = users.Filter(query)
	}

	usersJson, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		logger.Printf("%s", err)
		return createJsonError(ctx, 500, "Error marshalling users.")
	}

	return string(usersJson)
}

// Get a user.
func GetUser(ctx *web.Context, id string) string {
	setJsonHeaders(ctx)

	user, err := store.Load().GetUser(id)
	if err != nil {
		logger.Printf("%s", err)
		return createJsonError(ctx, 404, "Error finding user.")
	}

	userJson, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		logger.Printf("%s", err)
		return createJsonError(ctx, 500, "Error unmarshalling user.")
	}

	return string(userJson)
}

// Create a user.
func SaveUser(ctx *web.Context, id string) string {
	setJsonHeaders(ctx)

	// TODO(aaron): Error validation on the id.

	if err := store.Load().CreateUser(id); err != nil {
		logger.Printf("%s", err)
		return createJsonError(ctx, 500, "Error creating user.")
	}

	return "{}"
}

// Delete a user.
func DeleteUser(ctx *web.Context, id string) string {
	setJsonHeaders(ctx)

	return "OK"
}

// Create a new post.
func CreatePost(ctx *web.Context) string {
	setJsonHeaders(ctx)

	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		logger.Printf("%s", err)
		return createJsonError(ctx, 500, "Error reading request body.")
	}

	post := &store.Post{}
	if err := json.Unmarshal(body, post); err != nil {
		logger.Printf("%s", err)
		return createJsonError(ctx, 500, "Error unmarshalling post.")
	}

	if err := store.Load().CreatePost(post); err != nil {
		logger.Printf("%s", err)
		return createJsonError(ctx, 500, "Error saving post.")
	}

	postJson, err := json.MarshalIndent(post, "", "  ")
	if err != nil {
		logger.Printf("%s", err)
		return createJsonError(ctx, 500, "Error marshalling post.")
	}

	return string(postJson)
}

// Get a post.
func GetPost(ctx *web.Context, id string) string {
	setJsonHeaders(ctx)

	post, err := store.Load().GetPost(id)
	if err != nil {
		logger.Printf("%s", err)
		return createJsonError(ctx, 404, "Error finding post.")
	}

	postJson, err := json.MarshalIndent(post, "", "  ")
	if err != nil {
		logger.Printf("%s", err)
		return createJsonError(ctx, 500, "Error unmarshalling user.")
	}

	return string(postJson)
}

// Delete a post.
func DeletePost(ctx *web.Context, id string) string {
	setJsonHeaders(ctx)

	user, err := store.Load().GetVoleUser()
	if err != nil {
		logger.Printf("%s", err)
		return createJsonError(ctx, 500, "Error loading user.")
	}

	posts, err := store.Load().GetPosts(user.Id, "", config.ReadInt("UI_PageSize"))
	if err != nil {
		logger.Printf("%s", err)
		return createJsonError(ctx, 500, "Error loading posts.")
	}

	for _, post := range *posts {
		if post.Id == id {
			err := store.Load().DeletePost(&post)
			if err != nil {
				logger.Printf("%s", err)
				return createJsonError(ctx, 500, "Error deleting post.")
			} else {
				return ""
			}
		}
	}

	return ""
}

// Get a draft.
func GetDraft(ctx *web.Context, id string) string {
	setJsonHeaders(ctx)

	post := store.Post{}

	data, err := ioutil.ReadFile(path.Join(config.StorageDir(), "drafts", id+".json"))
	if err != nil {
		return createJsonError(ctx, 404, "Error loading draft.")
	}

	if err := json.Unmarshal(data, &post); err != nil {
		return createJsonError(ctx, 500, "No post or invalid json.")
	}

	rawJson, err := json.Marshal(post)
	return string(rawJson)
}

// Create a new draft.
func CreateDraft(ctx *web.Context) string {
	setJsonHeaders(ctx)

	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		logger.Printf("%s", err)
		return createJsonError(ctx, 500, "Error reading request body.")
	}

	draft := &store.Post{}
	if err := json.Unmarshal(body, draft); err != nil {
		logger.Printf("%s", err)
		return createJsonError(ctx, 500, "Error marshalling draft.")
	}

	// Create a new UUID for this draft.
	uuidBytes, _ := uuid.NewV4()
	draft.Id = fmt.Sprintf("%s", uuidBytes)

	// Drafts are just posts that aren't shared with other users.
	draftPath := path.Join(config.StorageDir(), "drafts", draft.Id+".json")

	draft.Created = time.Now().UnixNano()
	draft.Modified = time.Now().UnixNano()

	rawJson, err := json.MarshalIndent(draft, "", "  ")
	if err != nil {
		logger.Printf("%s", err)
		return createJsonError(ctx, 500, "Error marshalling draft.")
	}

	if err := store.Write(draftPath, rawJson); err != nil {
		logger.Printf("%s", err)
		return createJsonError(ctx, 500, "Error saving draft.")
	}

	return string(rawJson)
}

// Save a draft.
func SaveDraft(ctx *web.Context, id string) string {
	setJsonHeaders(ctx)

	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return createJsonError(ctx, 500, "Error reading request body.")
	}

	draft := &store.Post{}
	if err := json.Unmarshal(body, draft); err != nil {
		return createJsonError(ctx, 500, "Invalid JSON.")
	}

	draftPath := path.Join(config.StorageDir(), "drafts", fmt.Sprintf("%s.json", draft.Id))

	draft.Modified = time.Now().UnixNano()

	rawJson, err := json.MarshalIndent(draft, "", "  ")
	if err != nil {
		return createJsonError(ctx, 500, "Error saving draft.")
	}

	store.Write(draftPath, rawJson)

	return string(rawJson)
}

// Delete a draft.
func DeleteDraft(ctx *web.Context, id string) string {
	fullPath := path.Join(config.StorageDir(), "drafts", id+".json")

	if err := store.Delete(fullPath); err != nil {
		logger.Printf("%s", err)
		return createJsonError(ctx, 500, "Error deleting draft.")
	}

	return ""
}

// Get the full list of the user's drafts.
func GetDrafts(ctx *web.Context) string {
	setJsonHeaders(ctx)

	collection := make(store.PostCollection, 0)

	draftsPath := path.Join(config.StorageDir(), "drafts")
	postFiles, _ := ioutil.ReadDir(draftsPath)

	for _, postFile := range postFiles {
		if !strings.HasSuffix(postFile.Name(), ".json") {
			logger.Printf("Unexpected file: %s", postFile.Name())
			continue
		}

		fullPath := path.Join(draftsPath, postFile.Name())
		data, err := ioutil.ReadFile(fullPath)
		if err != nil {
			logger.Printf("%s", err)
			continue
		}

		post := store.Post{}
		if err := json.Unmarshal(data, &post); err != nil {
			logger.Printf("%s", err)
			return createJsonError(ctx, 500, "Error unmarshalling post.")
		}

		collection = append(collection, post)
	}

	sort.Sort(collection)

	rawJson, _ := json.Marshal(collection)
	return string(rawJson)
}
