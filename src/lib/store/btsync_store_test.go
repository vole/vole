package store

import (
	"os"
	"testing"
)

// TODO(aaron): Make this a flag?
var testDir string = "/tmp/VoleTest"

var store Store = NewBTSyncStore(testDir)

/**
 * Setup the test directory
 */
func TestSetup(t *testing.T) {
	os.RemoveAll(testDir)
	os.MkdirAll(testDir, 0755)
}

func TestGetAndSetVoleUser(t *testing.T) {
	saveUser := &User{}
	saveUser.Name = "Professor Vole"
	err := store.SaveVoleUser(saveUser)

	if err != nil {
		t.Error(err)
	}
	if len(saveUser.Id) != 36 {
		t.Error("UUID wasn't set to a 36 character string")
	}
	if saveUser.Name != "Professor Vole" {
		t.Error("User name wasn't set correctly")
	}

	getUser, err := store.GetVoleUser()
	if err != nil {
		t.Error("Unable to get Vole user")
	}
	if getUser.Id != saveUser.Id {
		t.Error("Get should have returned the previously created user")
	}
}

func TestGetAndSaveUsers(t *testing.T) {
	jim := &User{}
	jim.Name = "Jim"

	sally := &User{}
	sally.Name = "Sally"

	if err := store.SaveUser(jim); err != nil {
		t.Error("Error creating user: Jim")
	}

	if err := store.SaveUser(sally); err != nil {
		t.Error("Error creating user: Sally")
	}

	users, err := store.GetUsers()

	if err != nil {
		t.Error("Error getting users")
	}
	if len(*users) < 2 {
		t.Error("Expected at least 2 users but found none")
	}
	if users.FindById(jim.Id) == -1 {
		t.Error("Expected to find Jim")
	}
	if users.FindById(sally.Id) == -1 {
		t.Error("Expected to find Sally")
	}
}

func TestGetAndSavePost(t *testing.T) {
	user, err := store.GetVoleUser()
	if err != nil {
		t.Error("Error loading Vole user")
	}

	post := &Post{}
	post.Body = "test post please ignore"
	post.User = user

	if err := store.SavePost(post); err != nil {
		t.Error("Error saving post")
	}

	posts, err := store.GetPostsForUser(user)
	if err != nil {
		t.Error("Error getting Vole user's posts")
	}
	if posts == nil || posts.Len() != 1 {
		t.Error("Expected 1 post")
	}

	savedPost, err := store.GetPost(post.Id)
	if err != nil {
		t.Error("Could not find saved post")
	}
	if savedPost != nil && savedPost.Id != post.Id {
		t.Error("Expected post ids to match")
	}
	if savedPost != nil && savedPost.User.Id != user.Id {
		t.Error("Expected post user ids to match")
	}
}

func TestGetPosts(t *testing.T) {
	kelly := &User{}
	kelly.Name = "Kelly"

	reggie := &User{}
	reggie.Name = "Reggie"

	if err := store.SaveUser(kelly); err != nil {
		t.Error("Error saving Kelly", err)
	}
	if err := store.SaveUser(reggie); err != nil {
		t.Error("Error saving Reggie", err)
	}

	post1 := &Post{}
	post1.Body = "Kelly's post"
	post1.User = kelly

	post2 := &Post{}
	post2.Body = "Reggie's post"
	post2.User = reggie

	if err := store.SavePost(post1); err != nil {
		t.Error("Error saving Kelly's post", err)
	}
	if err := store.SavePost(post2); err != nil {
		t.Error("Error saving Reggie's post", err)
	}

	posts, err := store.GetPosts()
	if err != nil {
		t.Error("Error getting posts", err)
	}
	if posts != nil && posts.Len() < 2 {
		t.Error("Expected at least 2 posts")
	}
}
