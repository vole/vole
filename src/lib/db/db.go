package db

import (
  "encoding/json"
  "fmt"
  "github.com/vole/gouuid"
  "io/ioutil"
  "os"
  osuser "os/user"
  "path"
  "strings"
  "time"
  "errors"
)

const VERSION = "v1"

var DIR = func() string {
  dir := "."
  user, err := osuser.Current()

  if err == nil {
    dir = user.HomeDir
  }

  return path.Join(dir, "Vole")
}()

func ReadFile(args ...string) ([]byte, error) {
  return ioutil.ReadFile(path.Join(args...))
}

func ReadDir(args ...string) ([]os.FileInfo, error) {
  return ioutil.ReadDir(path.Join(args...))
}

func Create(args ...string) (*os.File, error) {
  return os.Create(path.Join(args...))
}

/**
 * Posts
 */

type Post struct {
  Id    string `json:"id"`
  Title   string `json:"title"`
  User    string `json:"user"`
  Created int64  `json:"created"`
}

type PostContainer struct {
  Post Post `json:"post"`
}

type PostCollection struct {
  Posts []Post `json:"posts"`
}

func (collection *PostCollection) Len() int {
  return len(collection.Posts)
}

func (collection *PostCollection) Less(i, j int) bool {
  return collection.Posts[i].Created < collection.Posts[j].Created
}

func (collection *PostCollection) Swap(i, j int) {
  collection.Posts[i], collection.Posts[j] = collection.Posts[j], collection.Posts[i]
}

func (post *Post) Save() error {
  if post.Id == "" {
    uuid, _ := uuid.NewV4()
    post.Id = fmt.Sprintf("%s", uuid)
  }

  if post.Created == 0 {
    post.Created = time.Now().UnixNano()
  }

  rawJson, err := json.Marshal(*post)
  if err != nil {
    return err
  }

  file, err := Create(DIR, "users", post.User, VERSION, "posts", post.FileName())
  if err != nil {
    return err
  }

  file.Write(rawJson)
  return nil
}

func (post *Post) FileName() string {
  return fmt.Sprintf("%d-post-%s", post.Created, post.Id)
}

func PostFromJson(rawJson []byte) (*Post, error) {
  var post Post
  json.Unmarshal(rawJson, &post)

  if post.Id == "" {
    return nil, errors.New("Unable to use this file")
  }

  return &post, nil
}

func PostContainerFromJson(rawJson []byte) *PostContainer {
  var container PostContainer
  json.Unmarshal(rawJson, &container)
  return &container
}

func NewPost(title string, user string) *Post {
  return &Post{
    Title: title,
    User:  user,
  }
}

func NewPostContainer(post *Post) *PostContainer {
  return &PostContainer{Post: *post}
}

func NewPostCollection(posts []Post) *PostCollection {
  return &PostCollection{Posts: posts}
}

func GetPosts() (*PostCollection, error) {
  posts := make([]Post, 0)

  users, _ := ReadDir(DIR, "users")

  for _, user := range users {
    dir := path.Join(DIR, "users", user.Name(), VERSION, "posts")

    files, err := ReadDir(dir)
    if err != nil {
      continue
    }

    for _, post := range files {
      data, err := ReadFile(dir, post.Name())

      if err != nil {
        continue
      }

      postData, err := PostFromJson(data)

      if err != nil {
        continue
      }

      posts = append(posts, *postData)
    }
  }

  return &PostCollection{posts}, nil
}

/**
 * Users
 */

type User struct {
  Key         string `json:"key"`
  Hash        string `json:"hash"`
  User        string `json:"user"`
  DisplayName string `json:"display_name"`
  IsMyUser    bool `json:"is_my_user"`
}

type UserCollection struct {
  Users []User `json:"users"`
}

type UserContainer struct {
  User User `json:"user"`
}

func UserFromJson(rawJson []byte) *User {
  var user User
  json.Unmarshal(rawJson, &user)
  user.IsMyUser = true
  return &user
}

func NewUser(user string, displayName string) *User {
  return &User{
    User:        user,
    DisplayName: displayName,
  }
}

func NewUserContainer(user *User) *UserContainer {
  return &UserContainer{User: *user}
}

func NewUserCollection(users []User) *UserCollection {
  return &UserCollection{Users: users}
}

func CurrentUser() (*User, error) {
  rawName, err := ReadFile(DIR, "my_user")
  if err != nil {
    return nil, err
  }

  name := strings.TrimSpace(string(rawName))

  data, err := ReadFile(DIR, "users", name, VERSION, "user", name)
  if err != nil {
    return nil, err
  }

  return UserFromJson(data), nil
}
