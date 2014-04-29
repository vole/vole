package api

import (
	"encoding/json"
	"fmt"
	btsync "github.com/vole/btsync-api"
	"github.com/vole/web"
	"io/ioutil"
	"lib/config"
	"lib/store"
	"os"
	"path"
)

type Error struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

var userStore = &store.UserStore{
	Path:    config.StorageDir(),
	Version: config.Version(),
}

func setJsonHeaders(ctx *web.Context) {
	ctx.ContentType("json")
	ctx.SetHeader("Cache-Control", "no-cache, no-store", true)
}

func createJsonError(ctx *web.Context, message string) string {
	ctx.ResponseWriter.WriteHeader(500)
	var err = Error{true, message}
	bytes, _ := json.Marshal(err)
	return string(bytes)
}

func Status(ctx *web.Context) string {
	setJsonHeaders(ctx)

	// TODO: Load from config.
	api := btsync.New("vole", "vole", 8888, true)

	_, err := api.GetVersion()
	if err != nil {
		return "{ \"btsync\": false }"
	}

	return "{ \"btsync\": true }"
}

/**
 * GetConfig()
 *
 * Fetch the app config.
 */
func GetConfig(ctx *web.Context) string {
	setJsonHeaders(ctx)

	configJson, err := config.Json()
	if err != nil {
		return createJsonError(ctx, "Error loading config.")
	}

	return configJson
}

/**
 * GetDrafts()
 *
 * Get all of the user's drafts.
 */
func GetDrafts(ctx *web.Context) string {
	setJsonHeaders(ctx)

	myUser, _ := userStore.GetMyUser()

	allDrafts, err := myUser.GetDrafts()
	if err != nil {
		return createJsonError(ctx, "Error loading drafts.")
	}

	draftsJson, err := allDrafts.Json()
	if err != nil {
		return createJsonError(ctx, "Error getting posts as json.")
	}

	return draftsJson
}

/**
 * GetPosts()
 */
func GetPosts(ctx *web.Context) string {
	setJsonHeaders(ctx)

	limit := config.ReadInt("UI_PageSize")
	before, _ := ctx.Params["before"]
	userId, _ := ctx.Params["user"]

	var allPosts *store.PostCollection
	var err error

	if userId != "" {
		var user *store.User
		if userId == "my_user" {
			user, err = userStore.GetMyUser()
		} else {
			user, err = userStore.GetUserById(userId)
		}
		if err != nil {
			return createJsonError(ctx, "User not found while getting posts.")
		}
		allPosts, err = user.GetPosts()
		if err != nil {
			return createJsonError(ctx, "Error loading posts.")
		}
	} else {
		allPosts, err = userStore.GetPosts()
		if err != nil || len(allPosts.Posts) < 1 {
			// Return a welcome post.
			post := &store.Post{}
			post.InitNew("Welcome to Vole. To start, create a new profile by clicking 'My Profile' on the left.", "none", "none", "Welcome", "", false)
			post.Id = "none"
			allPosts = post.Collection()
		}
	}

	allPosts.BeforeId(before)
	allPosts.Limit(limit)

	postsJson, err := allPosts.Json()
	if err != nil {
		return createJsonError(ctx, "Error getting posts as json.")
	}

	return postsJson
}

/**
 * GetUsers()
 */
func GetUsers(ctx *web.Context) string {
	setJsonHeaders(ctx)

	_, isMyUserFilter := ctx.Params["is_my_user"]
	query, hasQuery := ctx.Params["query"]

	var users *store.UserCollection
	var err error

	if isMyUserFilter {
		myUser, _ := userStore.GetMyUser()
		if myUser != nil {
			users = myUser.Collection()
		} else {
			users = store.GetEmptyUserCollection()
		}
	} else {
		users, err = userStore.GetUsers()
		if hasQuery {
			users.Filter(query)
		}
		if err != nil {
			return createJsonError(ctx, "Error loading all users.")
		}
	}

	usersJson, err := users.Json()
	if err != nil {
		return createJsonError(ctx, "Error getting users as json.")
	}

	return usersJson
}

/**
 * SaveUser()
 */
func SaveUser(ctx *web.Context) string {
	setJsonHeaders(ctx)

	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return createJsonError(ctx, "Error reading request body.")
	}

	user, err := userStore.NewUserFromContainerJson(body)
	if err != nil {
		return createJsonError(ctx, "Invalid JSON")
	}

	if err := user.Save(); err != nil {
		return createJsonError(ctx, "Error saving user")
	}

	if err := userStore.SetMyUser(user); err != nil {
		return createJsonError(ctx, "Error setting my user")
	}

	container := user.Container()

	userJson, err := container.Json()
	if err != nil {
		return createJsonError(ctx, "Could not create container")
	}

	return userJson
}

/**
 * SavePost()
 */
func SavePost(ctx *web.Context) string {
	setJsonHeaders(ctx)

	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return createJsonError(ctx, "Error reading request body.")
	}

	user, err := userStore.GetMyUser()
	if err != nil {
		return createJsonError(ctx, "Error reading my user when posting.")
	}

	post, err := user.NewPostFromJson(body)
	if err != nil {
		return createJsonError(ctx, "Invalid JSON")
	}

	if err := post.Save(); err != nil {
		return createJsonError(ctx, "Error saving post")
	}

	container := post.Container()

	postJson, err := container.Json()
	if err != nil {
		return createJsonError(ctx, "Could not create container")
	}

	return postJson
}

/**
 * GetPost()
 */
func GetPost(ctx *web.Context, id string) string {
	setJsonHeaders(ctx)

	user, err := userStore.GetMyUser()
	if err != nil {
		return createJsonError(ctx, "Error loading user.")
	}

	posts, err := user.GetPosts()
	if err != nil {
		return createJsonError(ctx, "Error loading posts.")
	}

	for _, post := range posts.Posts {
		if post.Id == id {
			rawJson, err := json.Marshal(post)
			if err != nil {
				return createJsonError(ctx, "Error parsing JSON.")
			}
			return string(rawJson)
		}
	}

	return createJsonError(ctx, "Error finding post.")
}

/**
 * GetDraft()
 */
func GetDraft(ctx *web.Context, id string) string {
	setJsonHeaders(ctx)

	user, err := userStore.GetMyUser()
	if err != nil {
		return createJsonError(ctx, "Error loading user.")
	}

	posts, err := user.GetDrafts()
	if err != nil {
		return createJsonError(ctx, "Error loading posts.")
	}

	for _, post := range posts.Posts {
		if post.Id == id {
			rawJson, err := json.Marshal(post)
			if err != nil {
				return createJsonError(ctx, "Error parsing JSON.")
			}
			return string(rawJson)
		}
	}

	return createJsonError(ctx, "Error finding post.")
}

/**
 * DeletePost()
 */
func DeletePost(ctx *web.Context, id string) string {
	setJsonHeaders(ctx)

	user, err := userStore.GetMyUser()
	if err != nil {
		return createJsonError(ctx, "Error loading user.")
	}

	posts, err := user.GetPosts()
	if err != nil {
		return createJsonError(ctx, "Error loading posts.")
	}

	for _, post := range posts.Posts {
		if post.Id == id {
			err := post.Delete()
			if err != nil {
				return createJsonError(ctx, "Error deleting post.")
			} else {
				return "OK"
			}
		}
	}

	return "OK"
}

/**
 * SaveFriend()
 */
func SaveFriend(ctx *web.Context) string {
	setJsonHeaders(ctx)

	key, _ := ctx.Params["key"]

	if key == "" {
		return createJsonError(ctx, "Key is empty")
	}

	os.Mkdir(path.Join(config.StorageDir(), "users", key), 0700)

	// TODO: Load from config.
	api := btsync.New("vole", "vole", 8888, true)

	response, err := api.AddFolderWithSecret(path.Join(config.StorageDir(), "users", key), key)
	if err != nil {
		return createJsonError(ctx, fmt.Sprintf("add_folder: %s", err))
	}

	responseJson, err := json.Marshal(response)
	return string(responseJson)
}
