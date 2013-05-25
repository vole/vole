package store

import (
  "encoding/json"
  "time"
  "path"
  "fmt"
  "github.com/vole/gouuid"
)

/**
 * Post.
 */
type Post struct {
  // Properties that should be saved to disk.
  Id           string `json:"id"`
  Title        string `json:"title"`
  Created      int64  `json:"created"`

  // Properties that are used by Vole backend and frontend, but not saved to disk
  // when the post is marshaled.
  UserId       string `json:"user_id,omitempty"`
  UserName     string `json:"user_name,omitempty"`
  UserAvatar string `json:"user_avatar,omitempty"`

  // Properties that are only used by the backend and thus don't have
  // to be marshaled to JSON for either the frontend or disk.
  FullPath     string `json:"-"`
}

type PostCollection struct {
  Posts []Post `json:"posts"`
}

func (collection *PostCollection) Len() int {
  return len(collection.Posts)
}

func (collection *PostCollection) Less(i, j int) bool {
  return collection.Posts[i].Created > collection.Posts[j].Created
}

func (collection *PostCollection) Swap(i, j int) {
  collection.Posts[i], collection.Posts[j] = collection.Posts[j], collection.Posts[i]
}

/**
 * InitNew()
 *
 * Initialize a new post creating the id and other fields.
 */
func (post *Post) InitNew(title, userPath, userId, userName, userAvatar string) {
  // Create a new UUID
  uuidBytes, _ := uuid.NewV4()
  uuid := fmt.Sprintf("%s", uuidBytes)

  // Get the timestamp.
  created := time.Now().UnixNano()

  // The full path to the post.
  fullPath := path.Join(userPath, "posts", fmt.Sprintf("%d-post-%s.json", created, uuid))

  post.Id = uuid
  post.Title = title
  post.Created = created
  post.UserId = userId
  post.UserName = userName
  post.UserAvatar = userAvatar
  post.FullPath = fullPath
}

/**
 * InitFromJson()
 *
 * Initialize a new post from json data
 */
func (post *Post) InitFromJson(rawJson []byte, fullPath string, userId string, userName string, userAvatar string) error {
  if err := json.Unmarshal(rawJson, post); err != nil {
    return err
  }
  post.UserId = userId
  post.UserName = userName
  post.UserAvatar = userAvatar
  post.FullPath = fullPath
  return nil
}

/**
 * Save()
 *
 * Save post to disk.
 */
func (post *Post) Save() error {
  // Before marshaling JSON for saving to disk, we set all properties
  // that should not be saved to empty, so they are ignored by marshaller.
  post.UserId = ""
  post.UserName = ""
  post.UserAvatar = ""

  rawJson, err := json.Marshal(*post)
  if err != nil {
    return err
  }

  return Write(post.FullPath, rawJson)
}
