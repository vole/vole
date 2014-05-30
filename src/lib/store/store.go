package store

import (
	"lib/config"
	"log"
	"os"
)

var (
	logger    = log.New(os.Stdout, "[Vole] ", log.Ldate|log.Ltime)
	dataStore Store
)

func Load() Store {
	if dataStore == nil {
		dataStore = NewBTSyncStore(config.StorageDir())
	}

	return dataStore
}

type StoreConfig struct{}

type Store interface {
	// Returns the name of the interface.
	Name() string
	// Version number for store.
	Version() string

	GetConfig() (*StoreConfig, error)
	SetConfig(*StoreConfig) error

	// Initialize the store.
	Initialize() error

	// Gets the current Vole user.
	GetVoleUser() (*User, error)
	// Creates a new Vole user.
	CreateVoleUser(user *User) error

	// Get a user.
	GetUser(id string) (*User, error)
	// Create a new user.
	CreateUser(id string) error
	// Get a list of all users.
	GetUsers() (*UserCollection, error)

	// Get a post.
	GetPost(id string) (*Post, error)
	// Create a post.
	CreatePost(post *Post) error
	// Delete a post.
	DeletePost(post *Post) error
	// Get a list of all posts.
	GetPosts(userId string, postWatermark string, limit int) (*PostCollection, error)
}
