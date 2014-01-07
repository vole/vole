package store

import (
  "os"
  osuser "os/user"
  "path"
  "strings"
  "testing"
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
  Path:    DIR,
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
  if !strings.Contains(user.FullPath, user.Name+"_"+user.Id) {
    t.Error("User path location not set correctly")
  }
  if !strings.Contains(user.UserJsonPath, user.Name+"_"+user.Id) {
    t.Error("User json path location not set correctly")
  }
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
  if user.IsMyUser != true {
    t.Error("After saving a new user, ismyuser must be true")
  }
  assertFileOrDir(t, path.Join(userStore.Path, "users", user.Name+"_"+user.Id, userStore.Version))
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

func TestGetAllPosts(t *testing.T) {
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

func TestPostsLimit(t *testing.T) {
  allPosts, err := userStore.GetPosts()
  if err != nil {
    t.Error(err)
  }

  allPosts.Limit(100)
  if len(allPosts.Posts) != 6 {
    t.Error("Large limit failure")
  }
  allPosts.Limit(0)
  if len(allPosts.Posts) != 6 {
    t.Error("Zero limit failure")
  }

  allPosts.Limit(2)
  if len(allPosts.Posts) != 2 {
    t.Error("Only two posts should be returned.")
  }
  for i := 0; i < 2; i++ {
    if allPosts.Posts[i].Title != titles[5-i] {
      t.Error("Title of post found to be incorrect, expected: '" + titles[5-i] + "' found: '" + allPosts.Posts[i].Title + "'")
    }
  }
}

func TestPostsBeforeId(t *testing.T) {
  allPosts, err := userStore.GetPosts()
  if err != nil {
    t.Error(err)
  }

  allPosts.BeforeId("garbage")
  if len(allPosts.Posts) != 6 {
    t.Error("Garbage ID failure")
  }

  beforeId := allPosts.Posts[0].Id
  allPosts.BeforeId(beforeId)
  if len(allPosts.Posts) != 5 {
    t.Error("Boundary limit fail")
  }

  allPosts, _ = userStore.GetPosts()
  beforeId = allPosts.Posts[2].Id
  allPosts.BeforeId(beforeId)
  if len(allPosts.Posts) != 3 {
    t.Error("Before ID failure")
  }
  for i := 0; i < 2; i++ {
    if allPosts.Posts[i].Title != titles[2-i] {
      t.Error("Title of post found to be incorrect, expected: '" + titles[2-i] + "' found: '" + allPosts.Posts[i].Title + "'")
    }
  }
}

func TestSaveUserFromJson(t *testing.T) {
  rawJson := []byte(`{"user":{"name":"json_user","avatar":null,"is_my_user":true,"email":"json_user@example.com"}}`)
  user, err := userStore.NewUserFromContainerJson(rawJson)
  if err != nil {
    t.Error(err)
  }

  if len(user.Id) != 36 {
    t.Error("UUID wasn't set to a 36 character string")
  }
  if user.Name != "json_user" {
    t.Error("Name was not 'json_user'")
  }
  if user.Email != "json_user@example.com" {
    t.Error("Email was not json_user@example.com'")
  }

  if err := user.Save(); err != nil {
    t.Error(err)
  }
  assertFileOrDir(t, path.Join(userStore.Path, "users", user.Name+"_"+user.Id, userStore.Version))

  container := user.Container()
  userJson, err := container.Json()
  if err != nil {
    t.Error("Unable to access container json")
  }
  if len(userJson) < 10 {
    t.Error("Json string suspiciously short")
  }
}

func TestCreatePostFromJson(t *testing.T) {
  user, _ := userStore.GetMyUser()
  rawJson := []byte(`{"post":{"title":"This is a test post","created":null,"user_id":null,"user_name":null,"user_avatar":null}}`)

  post, err := user.NewPostFromContainerJson(rawJson)
  if err != nil {
    t.Error(err)
  }

  if err := post.Save(); err != nil {
    t.Logf("%+v", post)
    t.Error(err)
  }
}
