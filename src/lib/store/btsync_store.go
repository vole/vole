package store

import (
	"code.google.com/p/go.exp/fsnotify"
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	btsync "github.com/vole/btsync-api"
	"github.com/vole/gouuid"
	"lib/config"
	"log"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
	"time"
)

var logger = log.New(os.Stdout, "[Vole] ", log.Ldate|log.Ltime)

// Create a new BTSyncStore.
func NewBTSyncStore(path string) Store {
	store := new(BTSyncStore)

	// Path to the store's directory.
	store.Path = path

	// Index of store's files.
	// TODO(aaron): Handle errors creating the index.
	store.Index = new(BTSyncStoreIndex)

	// Client for connecting to the BT Sync API.
	store.Client = btsync.New(
		config.ReadString("BTSync_User"),
		config.ReadString("BTSync_Pass"),
		config.ReadInt("BTSync_Port"),
		true)

	// Initialize the store.
	store.Initialize()

	return store
}

// An index of the file system
//
// TODO(aaron): The fsnotify approach is problematic with a large
// number of files. The file system's file descriptor limit can
// easily be reached, especially on OSX where the default is low.
// An option is to detect the file descriptor limit error, and fallback
// to polling.
// https://github.com/howeyc/fsnotify/issues/89#issuecomment-43451454
type BTSyncStoreIndex struct {
	Posts   *list.List
	Watcher *fsnotify.Watcher
}

// Start the indexer.
func (index *BTSyncStoreIndex) Start() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev.IsCreate() {
					index.Add(ev.Name)
				} else if ev.IsDelete() || ev.IsRename() {
					index.Remove(ev.Name)
				} else {
					logger.Print(ev)
				}
			case err := <-watcher.Error:
				logger.Println("Index Error:", err)
			}
		}
	}()

	index.Watcher = watcher

	return nil
}

// Adds a directory to the index's file system monitor.
func (index *BTSyncStoreIndex) Watch(dir string) error {
	if index.Watcher == nil {
		err := index.Start()
		if err != nil {
			logger.Print(err)
			return err
		}
	}

	err := index.Watcher.Watch(dir)
	if err != nil {
		logger.Print(err)
		return err
	}

	return nil
}

// Add a post to the index.
func (index *BTSyncStoreIndex) Add(file string) error {
	data, err := ReadFile(file)
	if err != nil {
		return err
	}

	post := Post{}
	if err := json.Unmarshal(data, &post); err != nil {
		return errors.New("No post or invalid json.")
	}

	for e := index.Posts.Front(); e != nil; e = e.Next() {
		if post.Created > e.Value.(Post).Created {
			index.Posts.InsertBefore(post, e)
			logger.Printf("Adding new post %s before %s", post.Id, e.Value.(Post).Id)
			break
		}
	}

	return nil
}

// Remove a post from the index.
func (index *BTSyncStoreIndex) Remove(file string) error {
	re := regexp.MustCompile("[0-9]+-post-(.*).json")
	match := re.FindStringSubmatch(file)

	if match == nil {
		return errors.New("Error parsing Id from file name: " + file)
	}

	id := match[1]

	for e := index.Posts.Front(); e != nil; e = e.Next() {
		if id > e.Value.(Post).Id {
			index.Posts.Remove(e)
			logger.Printf("Removing post %s from index", id)
			break
		}
	}

	return nil
}

// A Vole store for Bittorrent Sync.
type BTSyncStore struct {
	Path   string
	Client *btsync.BTSyncAPI
	Index  *BTSyncStoreIndex
	Store
}

// Name of the store.
func (store *BTSyncStore) Name() string {
	return "Bittorrent Sync Store"
}

// Version of the store.
func (store *BTSyncStore) Version() string {
	return "v1"
}

// Initialize the data store.
func (store *BTSyncStore) Initialize() error {
	// If BTSync_Watch is disabled, we dumbly rebuild the index
	// every few seconds.
	// TODO(aaron): Use a smarter polling mechanism.
	if !config.ReadBool("BTSync_Watch") {
		ticker := time.NewTicker(time.Minute)
		go func() {
			for _ = range ticker.C {
				store.BuildIndex()
			}
		}()
	} else {
		logger.Print("BTSync_Watch is enabled, using fsnotify!")
	}

	// If BTSync_Watch is enabled, we can depend on the index
	// to watch the file system.
	return store.BuildIndex()
}

// TODO(aaron): To save memory, instead of storing an entire post in the index,
// store a reference to the file.
func (store *BTSyncStore) BuildIndex() error {
	now := time.Now()

	// Ignore posts older than 6 months.
	// TODO(aaron): Make this configurable.
	cutoff := now.AddDate(0, -6, 0)

	logger.Printf("Building index starting from: %s", cutoff)

	users, err := store.GetUsers()
	if err != nil {
		logger.Printf("Error getting users for index: %s", err)
		return err
	}

	collection := make(PostCollection, 0)

	for _, user := range *users {
		userPath := path.Join(store.Path, "users", user.Id, store.Version())
		postFiles, _ := ReadDir(userPath, "posts")

		for _, postFile := range postFiles {
			if !strings.HasSuffix(postFile.Name(), ".json") {
				continue
			}

			fullPath := path.Join(userPath, "posts", postFile.Name())
			data, err := ReadFile(fullPath)
			if err != nil {
				logger.Printf("Error reading file: %s", err)
				continue
			}

			post := Post{}
			if err := json.Unmarshal(data, &post); err != nil {
				logger.Printf("No post or invalid json.")
				continue
			}

			// TODO(aaron): I can't even begin to tell you why
			// I have to do this. Needless to say, I'm embarrassed.
			userCopy := new(User)
			*userCopy = user
			post.User = userCopy

			if post.Created > cutoff.UnixNano() {
				collection = append(collection, post)
			}
		}

		if config.ReadBool("BTSync_Watch") {
			store.Index.Watch(path.Join(userPath, "posts"))
		}
	}

	elapsed := time.Since(now)
	logger.Printf("Added %d posts to the index in %s", collection.Len(), elapsed)

	sort.Sort(collection)

	store.Index.Posts = list.New()
	for _, post := range collection {
		store.Index.Posts.PushBack(post)
	}

	return nil
}

// Get the Vole user.
func (store *BTSyncStore) GetVoleUser() (*User, error) {
	id, err := ReadFile(store.Path, "my_user")
	if err != nil {
		return nil, err
	}

	user, err := store.GetUser(string(id))
	if err != nil {
		return nil, err
	}

	user.IsMyUser = true

	return user, nil
}

// Create a new Vole user.
func (store *BTSyncStore) CreateVoleUser(user *User) error {
	// Create new Sync keys.
	secrets, err := store.Client.GetSecrets(false)
	if err != nil {
		return err
	}

	// Use the read-only key as the user's ID.
	user.Id = secrets.ReadOnly

	// Create the user's directory and save their information to disk.
	userDir := path.Join(store.Path, "users", user.Id)
	userPath := path.Join(userDir, store.Version(), "user")
	postsPath := path.Join(userDir, store.Version(), "posts")
	filesPath := path.Join(userDir, store.Version(), "files")

	// Create the user's user directory.
	if err := os.MkdirAll(userPath, 0755); err != nil {
		return err
	}

	// Create the user's posts directory.
	if err := os.MkdirAll(postsPath, 0755); err != nil {
		return err
	}

	// Create the user's files directory.
	if err := os.MkdirAll(filesPath, 0755); err != nil {
		return err
	}

	// Generate human readable JSON for user.
	rawJson, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		return err
	}

	jsonPath := path.Join(userPath, "user.json")

	// Save the user's JSON data.
	if err := Write(jsonPath, rawJson); err != nil {
		return err
	}

	// Register the new user directory with Sync.
	if _, err := store.Client.AddFolderWithSecret(userDir, secrets.ReadWrite); err != nil {
		return err
	}

	// Write the user's ID to the my_user file.
	if err := Write(path.Join(store.Path, "my_user"), []byte(user.Id)); err != nil {
		return err
	}

	return nil
}

// Get a user by their Id.
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

// Create a new user with the given Id.
func (store *BTSyncStore) CreateUser(id string) error {
	userDir := path.Join(config.StorageDir(), "users", id)

	if err := os.Mkdir(userDir, 0700); err != nil {
		return err
	}

	response, err := store.Client.AddFolderWithSecret(userDir, id)
	if err != nil {
		// If we can't add the folder to BT Sync, attempt to clean
		// up the directory we just created.
		if err := os.RemoveAll(userDir); err != nil {
			return err
		}

		return err
	}

	if response.Error != 0 {
		return errors.New(response.Message)
	}

	return nil
}

// Delete the user with the given Id.
func (store *BTSyncStore) DeleteUser(id string) error {
	// Path to the user's directory.
	userDir := path.Join(config.StorageDir(), "users", id)

	// Delete the user's directory.
	if err := os.RemoveAll(userDir); err != nil {
		return err
	}

	// Remove the user's directory from BT Sync.
	if _, err := store.Client.RemoveFolder(id); err != nil {
		return err
	}

	return nil
}

// Get a list of all users.
func (store *BTSyncStore) GetUsers() (*UserCollection, error) {
	// Grab the Vole user so we can flag them in the response.
	voleUser, err := store.GetVoleUser()
	if err != nil {
		return nil, err
	}

	collection := make(UserCollection, 0)

	// List of all user directories.
	userDirs, _ := ReadDir(store.Path, "users")

	// Load the user for each user directory.
	for _, dir := range userDirs {
		user, err := store.GetUser(dir.Name())
		if err != nil {
			continue
		}

		// If this is the Vole user, flag them.
		if user.Id == voleUser.Id {
			user.IsMyUser = true
		}

		collection = append(collection, *user)
	}

	return &collection, nil
}

// Create a new post.
func (store *BTSyncStore) CreatePost(post *Post) error {
	user, err := store.GetVoleUser()
	if err != nil {
		return err
	}

	// Give the post a UUID.
	uuidBytes, _ := uuid.NewV4()
	post.Id = fmt.Sprintf("%s", uuidBytes)

	// Someday if we allow editing of posts, the Modified date
	// won't be useless.
	post.Created = time.Now().UnixNano()
	post.Modified = time.Now().UnixNano()
	post.User = user

	store.Index.Posts.PushFront(*post)

	// User's posts directory.
	postDir := path.Join(store.Path, "users", user.Id, store.Version(), "posts")

	// Path to save new post.
	postPath := path.Join(postDir, fmt.Sprintf("%d-post-%s.json", post.Created, post.Id))

	rawJson, err := json.MarshalIndent(post, "", "  ")
	if err != nil {
		return err
	}

	return Write(postPath, rawJson)
}

// Delete a post.
func (store *BTSyncStore) DeletePost(post *Post) error {
	for e := store.Index.Posts.Front(); e != nil; e = e.Next() {
		if e.Value.(Post).Id == post.Id {
			store.Index.Posts.Remove(e)
		}
	}

	// Path to the post's user's posts directory.
	postDir := path.Join(store.Path, "users", post.User.Id, store.Version(), "posts")

	// Path to this specific post.
	postPath := path.Join(postDir, fmt.Sprintf("%d-post-%s.json", post.Created, post.Id))

	return Delete(postPath)
}

// Get a post.
func (store *BTSyncStore) GetPost(id string) (*Post, error) {
	for e := store.Index.Posts.Front(); e != nil; e = e.Next() {
		if e.Value.(Post).Id == id {
			return e.Value.(*Post), nil
		}
	}

	return nil, errors.New("Post not found")
}

// Get all posts.
func (store *BTSyncStore) GetPosts(userId string, postWatermark string, limit int) (*PostCollection, error) {
	collection := make(PostCollection, 0)

	postWatermarkFound := false

	for e := store.Index.Posts.Front(); e != nil; e = e.Next() {
		post := e.Value.(Post)

		if userId != "" && post.User.Id != userId {
			continue
		}

		if postWatermark != "" {
			if post.Id == postWatermark {
				postWatermarkFound = true
				continue
			} else if !postWatermarkFound {
				continue
			}
		}

		collection = append(collection, post)

		if limit != -1 && collection.Len() == limit {
			break
		}
	}

	return &collection, nil
}
