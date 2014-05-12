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

var dataStore = store.Load()

func setJsonHeaders(ctx *web.Context) {
	ctx.ContentType("json")
	ctx.SetHeader("Cache-Control", "no-cache, no-store", true)
}

func createJsonError(ctx *web.Context, message string) string {
	var err = Error{true, message}
	bytes, _ := json.Marshal(err)
	return string(bytes)
}

func Status(ctx *web.Context) string {
	setJsonHeaders(ctx)

	// TODO: Load from config.
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
 * GetPosts()
 */
func GetPosts(ctx *web.Context) string {
	setJsonHeaders(ctx)

	limit := config.ReadInt("UI_PageSize")
	before, _ := ctx.Params["before"]
	userId, _ := ctx.Params["user"]

	var posts *store.PostCollection
	var err error

	if userId != "" {
		user, err := dataStore.GetUser(userId)
		if err != nil {
			return createJsonError(ctx, "User not found")
		}

		posts, err = dataStore.GetPostsForUser(user)
	} else {
		posts, err = dataStore.GetPosts()
	}

	if err != nil {
		return createJsonError(ctx, fmt.Sprintf("%s", err))
	}

	if before != "" {
		posts.BeforeId(before)
	}

	posts.Limit(limit)

	postsJson, err := json.MarshalIndent(posts, "", "  ")
	if err != nil {
		return createJsonError(ctx, "Error getting posts as json.")
	}

	return string(postsJson)
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
	var usersJson []byte

	if isMyUserFilter {
		user, _ := dataStore.GetVoleUser()

		usersJson, err = json.MarshalIndent(user, "", "  ")
		if err != nil {
			return createJsonError(ctx, "Error getting users as json.")
		}
	} else {
		users, err = dataStore.GetUsers()
		if hasQuery {
			users.Filter(query)
		}
		if err != nil {
			return createJsonError(ctx, "Error loading all users.")
		}

		usersJson, err = json.MarshalIndent(users, "", "  ")
		if err != nil {
			return createJsonError(ctx, "Error getting users as json.")
		}
	}

	return string(usersJson)
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

	fmt.Println(string(body))

	var user = &store.User{}
	if err := json.Unmarshal(body, user); err != nil {
		fmt.Println(err)
		return createJsonError(ctx, "Error unmarshalling user")
	}

	if err := dataStore.SaveVoleUser(user); err != nil {
		return createJsonError(ctx, "Error saving user")
	}

	userJson, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		return createJsonError(ctx, "Could not create collection")
	}

	return string(userJson)
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

	user, err := dataStore.GetVoleUser()
	if err != nil {
		return createJsonError(ctx, "Error reading my user when posting.")
	}

	post := &store.Post{}
	if err := json.Unmarshal(body, post); err != nil {
		return createJsonError(ctx, "Invalid JSON")
	}

	post.User = user

	if err := dataStore.SavePost(post); err != nil {
		return createJsonError(ctx, "Error saving post")
	}

	postJson, err := json.MarshalIndent(post, "", "  ")
	if err != nil {
		return createJsonError(ctx, "Could not create container")
	}

	return string(postJson)
}

/**
 * GetPost()
 */
func GetPost(ctx *web.Context, id string) string {
	setJsonHeaders(ctx)

	post, err := dataStore.GetPost(id)
	if err != nil {
		return createJsonError(ctx, "Error finding post")
	}

	postJson, err := json.MarshalIndent(post, "", "  ")
	if err != nil {
		return createJsonError(ctx, "Error unmarshalling json")
	}

	return string(postJson)
}

/**
 * DeletePost()
 */
func DeletePost(ctx *web.Context, id string) string {
	setJsonHeaders(ctx)

	user, err := dataStore.GetVoleUser()
	if err != nil {
		return createJsonError(ctx, "Error loading user.")
	}

	posts, err := dataStore.GetPostsForUser(user)
	if err != nil {
		return createJsonError(ctx, "Error loading posts.")
	}

	for _, post := range *posts {
		if post.Id == id {
			err := dataStore.DeletePost(&post)
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

func GetDraft(ctx *web.Context, id string) string {
	setJsonHeaders(ctx)

	post := store.Post{}

	data, err := ioutil.ReadFile(path.Join(config.StorageDir(), "drafts", id+".json"))
	if err != nil {
		return createJsonError(ctx, "Error reading draft.")
	}

	if err := json.Unmarshal(data, &post); err != nil {
		return createJsonError(ctx, "No post or invalid json.")
	}

	rawJson, err := json.Marshal(post)
	return string(rawJson)
}

func CreateDraft(ctx *web.Context) string {
	setJsonHeaders(ctx)

	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return createJsonError(ctx, "Error reading request body.")
	}

	draft := &store.Post{}
	if err := json.Unmarshal(body, draft); err != nil {
		return createJsonError(ctx, "Invalid JSON.")
	}

	uuidBytes, _ := uuid.NewV4()
	draft.Id = fmt.Sprintf("%s", uuidBytes)

	draftPath := path.Join(config.StorageDir(), "drafts", fmt.Sprintf("%s.json", draft.Id))

	draft.Created = time.Now().UnixNano()
	draft.Modified = time.Now().UnixNano()

	rawJson, err := json.MarshalIndent(draft, "", "  ")
	if err != nil {
		return createJsonError(ctx, "Error saving draft.")
	}

	store.Write(draftPath, rawJson)

	return string(rawJson)
}

func SaveDraft(ctx *web.Context, id string) string {
	setJsonHeaders(ctx)

	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return createJsonError(ctx, "Error reading request body.")
	}

	draft := &store.Post{}
	if err := json.Unmarshal(body, draft); err != nil {
		return createJsonError(ctx, "Invalid JSON.")
	}

	draftPath := path.Join(config.StorageDir(), "drafts", fmt.Sprintf("%s.json", draft.Id))

	draft.Modified = time.Now().UnixNano()

	rawJson, err := json.MarshalIndent(draft, "", "  ")
	if err != nil {
		return createJsonError(ctx, "Error saving draft.")
	}

	store.Write(draftPath, rawJson)

	return string(rawJson)
}

// TODO(aaron): Handle errors.
func DeleteDraft(ctx *web.Context, id string) string {
	fullPath := path.Join(config.StorageDir(), "drafts", id+".json")
	store.Delete(fullPath)
	return "OK"
}

func GetDrafts(ctx *web.Context) string {
	setJsonHeaders(ctx)

	collection := make(store.PostCollection, 0)

	draftsPath := path.Join(config.StorageDir(), "drafts")
	postFiles, _ := ioutil.ReadDir(draftsPath)

	for _, postFile := range postFiles {
		if !strings.HasSuffix(postFile.Name(), ".json") {
			continue
		}

		fullPath := path.Join(draftsPath, postFile.Name())
		data, err := ioutil.ReadFile(fullPath)
		if err != nil {
			continue
		}

		post := store.Post{}
		if err := json.Unmarshal(data, &post); err != nil {
			return createJsonError(ctx, "No post or invalid json.")
		}

		collection = append(collection, post)
	}

	sort.Sort(collection)

	rawJson, _ := json.Marshal(collection)
	return string(rawJson)
}
