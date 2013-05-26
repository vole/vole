/**
 * Parts of the user struct that provide access to the user's posts.
 */

package store

import (
  "strings"
  "path"
  "errors"
  "encoding/json"
)

/**
 * NewPostFromContainerJson()
 *
 * Called by POST requests from the frontend.
 */
func (user *User) NewPostFromContainerJson(rawJson []byte) (*Post, error) {
  var container PostContainer
  if err := json.Unmarshal(rawJson, &container); err != nil {
    return nil, err
  }
  post := user.NewPost(container.Post.Title)
  return post, nil
}

/**
 * NewPost(title)
 *
 * Return a new Post struct for this user.
 */
func (user *User) NewPost(title string) *Post {
  post := &Post{}
  post.InitNew(title, user.FullPath, user.Id, user.Name, user.Avatar)
  return post
}

/**
 * GetPosts()
 *
 * Get all a user's posts.
 */
func (user *User) GetPosts() (*PostCollection, error) {
  collection := make([]Post, 0)
  postFiles, _ := ReadDir(user.FullPath, "posts")

  for _, postFile := range postFiles {
    if !strings.HasSuffix(postFile.Name(), ".json") {
      continue
    }
    fullPath := path.Join(user.FullPath, "posts", postFile.Name())
    data, err := ReadFile(fullPath)
    if err != nil {
      continue
    }
    post := Post{}
    if err := post.InitFromJson(data, fullPath, user.Id, user.Name, user.Avatar); err != nil {
      return nil, errors.New("No post or invalid json.")
    }
    collection = append(collection, post)
  }
  return &PostCollection{collection}, nil
}
