package store

import (
	"lib/config"
)

func Load() Store {
	// TODO(aaron): Add others :P
	return NewBTSyncStore(config.StorageDir())
}

type StoreConfig struct {
}

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
	// Saves the current Vole user.
	SaveVoleUser(user *User) error

	// Get a user.
	GetUser(id string) (*User, error)
	// Save a user.
	SaveUser(user *User) error
	// Get a list of all users.
	GetUsers() (*UserCollection, error)

	// Get a post.
	GetPost(id string) (*Post, error)
	// Save a post.
	SavePost(post *Post) error
	// Delete a post.
	DeletePost(post *Post) error
	// Get a list of all posts.
	GetPosts() (*PostCollection, error)
	// Get a list of all posts for a user.
	GetPostsForUser(user *User) (*PostCollection, error)
	// Get a list of all posts before a specific post.
	GetPostsBefore(id string) (*PostCollection, error)
	// Get a list of all posts after a specific post.
	GetPostsAfter(id string) (*PostCollection, error)
}
