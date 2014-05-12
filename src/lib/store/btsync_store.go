package store

import (
	"encoding/json"
	"errors"
	"fmt"
	btsync "github.com/vole/btsync-api"
	"github.com/vole/gouuid"
	"lib/config"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
	"time"
)

func NewBTSyncStore(path string) Store {
	store := new(BTSyncStore)
	store.Path = path
	store.Client = btsync.New(
		config.ReadString("BTSync_User"),
		config.ReadString("BTSync_Pass"),
		config.ReadInt("BTSync_Port"),
		true)
	return store
}

type BTSyncStoreConfig struct {
	StoreConfig
}

type BTSyncStore struct {
	Path   string
	Client *btsync.BTSyncAPI
	Store
}

func (store *BTSyncStore) Name() string {
	return "Bittorrent Sync Store"
}

func (store *BTSyncStore) Version() string {
	return "v1"
}

func (store *BTSyncStore) Initialize() error {
	return nil
}

func (store *BTSyncStore) GetVoleUser() (*User, error) {
	id, err := ReadFile(store.Path, "my_user")
	if err != nil {
		return nil, err
	}

	return store.GetUser(string(id))
}

func (store *BTSyncStore) SaveVoleUser(user *User) error {
	if err := store.SaveUser(user); err != nil {
		return err
	}

	if err := Write(path.Join(store.Path, "my_user"), []byte(user.Id)); err != nil {
		return err
	}

	return nil
}

func (store *BTSyncStore) GetUser(id string) (*User, error) {
	data, err := ReadFile(store.Path, "users", id, store.Version(), "user", "user.json")
	if err != nil {
		return nil, err
	}

	user := &User{}
	if err := json.Unmarshal(data, user); err != nil {
		return nil, err
	}

	// Convert old UUIDs to BTSync IDs.
	var uuid = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

	if uuid.MatchString(user.Id) {
		user.Id = id
	}

	return user, nil
}

func (store *BTSyncStore) SaveUser(user *User) error {
	userDir := path.Join(store.Path, "users", user.Id, store.Version())
	userPath := path.Join(userDir, "user")
	postsPath := path.Join(userDir, "posts")

	// Create the user's user directory.
	if err := os.MkdirAll(userPath, 0755); err != nil {
		return err
	}

	// Create the user's posts directory.
	if err := os.MkdirAll(postsPath, 0755); err != nil {
		return err
	}

	// Generate human readable JSON for user.
	rawJson, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		return err
	}

	jsonPath := path.Join(userPath, "user.json")

	// Save the user's JSON data.
	return Write(jsonPath, rawJson)
}

func (store *BTSyncStore) GetUsers() (*UserCollection, error) {
	collection := make(UserCollection, 0)
	userDirs, _ := ReadDir(store.Path, "users")

	for _, dir := range userDirs {
		user, err := store.GetUser(dir.Name())
		if err != nil {
			continue
		}

		collection = append(collection, *user)
	}

	return &collection, nil
}

func (store *BTSyncStore) SavePost(post *Post) error {
	postDir := path.Join(store.Path, "users", post.User.Id, store.Version(), "posts")

	if post.Id == "" {
		uuidBytes, _ := uuid.NewV4()
		post.Id = fmt.Sprintf("%s", uuidBytes)
	}

	postPath := path.Join(postDir, fmt.Sprintf("%d-post-%s.json", post.Created, post.Id))

	if post.Created == 0 {
		post.Created = time.Now().UnixNano()
	}

	post.Modified = time.Now().UnixNano()

	rawJson, err := json.MarshalIndent(post, "", "  ")
	if err != nil {
		return err
	}

	return Write(postPath, rawJson)
}

func (store *BTSyncStore) DeletePost(post *Post) error {
	postDir := path.Join(store.Path, "users", post.User.Id, store.Version(), "posts")
	postPath := path.Join(postDir, fmt.Sprintf("%d-post-%s.json", post.Created, post.Id))
	return Delete(postPath)
}

func (store *BTSyncStore) GetPost(id string) (*Post, error) {
	posts, err := store.GetPosts()
	if err != nil {
		return nil, err
	}

	index := posts.FindById(id)
	if index != -1 {
		return &(*posts)[index], nil
	}

	return nil, errors.New("Post not found")
}

func (store *BTSyncStore) GetPostsForUser(user *User) (*PostCollection, error) {
	collection := make(PostCollection, 0)

	userPath := path.Join(store.Path, "users", user.Id, store.Version())
	postFiles, _ := ReadDir(userPath, "posts")

	for _, postFile := range postFiles {
		if !strings.HasSuffix(postFile.Name(), ".json") {
			continue
		}

		fullPath := path.Join(userPath, "posts", postFile.Name())
		data, err := ReadFile(fullPath)
		if err != nil {
			continue
		}

		post := Post{}
		if err := json.Unmarshal(data, &post); err != nil {
			return nil, errors.New("No post or invalid json.")
		}

		// TODO(aaron): I can't even begin to tell you why
		// I have to do this. Needless to say, I'm embarrassed.
		userCopy := new(User)
		*userCopy = *user
		post.User = userCopy

		collection = append(collection, post)
	}

	sort.Sort(collection)
	return &collection, nil
}

func (store *BTSyncStore) GetPosts() (*PostCollection, error) {
	collection := make(PostCollection, 0)

	users, _ := store.GetUsers()

	for _, user := range *users {
		posts, err := store.GetPostsForUser(&user)
		if err != nil {
			continue
		}

		collection = append(collection, *posts...)
	}

	sort.Sort(collection)

	return &collection, nil
}

func (store *BTSyncStore) GetPostsBefore(id string) (*PostCollection, error) {
	posts, err := store.GetPosts()
	if err != nil {
		return nil, err
	}

	filtered := make(PostCollection, 0)

	i := posts.FindById(id)
	if i == -1 {
		return &filtered, nil
	}

	start := i + 1
	if start < posts.Len() {
		filtered = (*posts)[start:]
	}

	return &filtered, nil
}

func (store *BTSyncStore) GetPostsAfter(id string) (*PostCollection, error) {
	posts, err := store.GetPosts()
	if err != nil {
		return nil, err
	}

	filtered := make(PostCollection, 0)

	i := posts.FindById(id)
	if i == -1 {
		return &filtered, nil
	}

	start := i - 1
	if start > 0 {
		filtered = (*posts)[0 : i-1]
	}

	return &filtered, nil
}
