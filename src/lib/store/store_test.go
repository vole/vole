package store

import (
  "testing"
  "strings"
  //"fmt"
  "os"
  osuser "os/user"
  "path"
)

var DIR = func() string {
  dir := "."
  user, err := osuser.Current()
  if err == nil {
    dir = user.HomeDir
  }

  // WARNING: This directory is completely deleted if it exists. Don't mistype.
  return path.Join(dir, "VoleTest")
}()

var userStore = &UserStore{
  Path: DIR,
  Version: "v1",
}

func assertFileOrDir(t *testing.T, myUserPath string) {
  _, err := os.Stat(myUserPath)
  if err != nil {
    t.Error(err)
  }
}

/**
 * Setup the test directory
 */
func TestSetup(t *testing.T) {
  os.RemoveAll(userStore.Path)
  os.MkdirAll(userStore.Path, 0755)
}

func TestCreateUser(t *testing.T) {
  user := userStore.NewUser("first_user", "test@example.com")

  if len(user.Id) != 36 {
    t.Error("UUID wasn't set to a 36 character string")
  }
  if user.Name != "first_user" {
    t.Error("User name wasn't set correctly")
  }
  if user.Email != "test@example.com" {
    t.Error("Email not set correctly")
  }
  if !user.IsMyUser {
    t.Error("IsMyUser not set to true")
  }
  if !strings.Contains(user.Avatar, "https://gravatar.com/") {
    t.Error("Gravatar not being set correctly")
  }
  if !strings.Contains(user.FullPath, user.Name + "_" + user.Id) {
    t.Error("User path location not set correctly")
  }
  if !strings.Contains(user.UserJsonPath, user.Name + "_" + user.Id) {
    t.Error("User json path location not set correctly")
  }
  //t.Log(user)
}

func TestSaveUser(t *testing.T) {
  user := userStore.NewUser("first_user", "test@example.com")

  if err := userStore.SetMyUser(user); err != nil {
    t.Error(err)
  }
  assertFileOrDir(t, path.Join(userStore.Path, "my_user"))

  if err := user.Save(); err != nil {
    t.Error(err)
  }
  assertFileOrDir(t, path.Join(userStore.Path, "users", user.Name + "_" + user.Id, userStore.Version))
}

func TestGetMyUser(t *testing.T) {
  // Test for the user we just created in the function above.
  user, err := userStore.GetMyUser()
  if err != nil {
    t.Error(err)
  }

  if len(user.Id) != 36 {
    t.Error("UUID wasn't set to a 36 character string")
  }
  if user.Name != "first_user" {
    t.Error("User name wasn't retreived correctly")
  }
  if !user.IsMyUser {
    t.Error("User isn't my user")
  }
}

func TestGetAllUsers(t *testing.T) {
  // Add a user and then get them all.
  user := userStore.NewUser("second_user", "second_test@example.com")
  userStore.SetMyUser(user)
  user.Save()

  users, err := userStore.GetUsers()
  if err != nil {
    t.Error(err)
  }

  if len(users.Users) != 2 {
    t.Error("Incorrect number of users returned")
  }
  if users.Users[0].Name != "first_user" {
    t.Error("First user not found")
  }
  if users.Users[1].Name != "second_user" {
    t.Error("Second user not found")
  }
  if users.Users[0].IsMyUser {
    t.Error("First user should not be my user")
  }
  if !users.Users[1].IsMyUser {
    t.Error("Second user should be my user")
  }
}

func TestCreatePost(t *testing.T) {
  user, err := userStore.GetMyUser()
  if err != nil {
    t.Error(err)
  }

  post := user.NewPost("this is a test post")
  if len(post.Id) != 36 {
    t.Error("UUID wasn't set to a 36 character string")
  }
  if len(post.UserId) != 36 {
    t.Error("User UUID wasn't set to a 36 character string")
  }
  if post.Title != "this is a test post" {
    t.Error("Post title incorrectly set")
  }
}

var titles = [6]string{
  "this is a test post",
  "this is a second post",
  "making a third post",
  "and a fourth post",
  "my name is first_user and I enjoy posting",
  "my name is second_user and I enjoy posting more than first_user",
}

func TestSavePost(t *testing.T) {
  user, err := userStore.GetMyUser()
  if err != nil {
    t.Error(err)
  }
  post := user.NewPost(titles[0])
  post.Save()
  assertFileOrDir(t, post.FullPath)
}

func TestGetUserPosts(t *testing.T) {
  user, err := userStore.GetMyUser()
  if err != nil {
    t.Error(err)
  }
  post := user.NewPost(titles[1])
  post.Save()
  post = user.NewPost(titles[2])
  post.Save()
  post = user.NewPost(titles[3])
  post.Save()

  posts, err := user.GetPosts()
  if err != nil {
    t.Error(err)
  }

  for i := 0; i < 4; i++ {
    if posts.Posts[i].Title != titles[i] {
      t.Error("Title of post found to be incorrect, expected: '" + titles[i] + "' found: '" + posts.Posts[i].Title + "'")
    }
  }
}

func TestGetAllUsersPosts(t *testing.T) {
  users, err := userStore.GetUsers()
  if err != nil {
    t.Error(err)
  }
  user := &users.Users[0]
  userStore.SetMyUser(user)
  post := user.NewPost(titles[4])
  post.Save()

  user = &users.Users[1]
  userStore.SetMyUser(user)
  post = user.NewPost(titles[5])
  post.Save()

  allPosts, err := userStore.GetPosts()
  if err != nil {
    t.Error(err)
  }

  // Reverse order of posting
  j := 5
  for _, post := range allPosts.Posts {
    if post.Title != titles[j] {
      t.Error("Title of post found to be incorrect, expected: '" + titles[j] + "' found: '" + post.Title + "'")
    }
    j--
  }
}









