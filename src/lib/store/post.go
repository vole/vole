package store

import (
  "encoding/json"
  "fmt"
  "github.com/vole/gouuid"
  "path"
  "time"
)

/**
 * Post.
 */
type Post struct {
  // Properties that should be saved to disk.
  Id      string `json:"id"`
  Title   string `json:"title"`
  Created int64  `json:"created"`
  RepostUserName   string `json:"repost_user_name,omitempty"`
  RepostUserAvatar string `json:"repost_user_avatar,omitempty"`
  RepostPostId     string `json:"repost_post_id,omitempty"`

  // Properties that are used by Vole backend and frontend, but not saved to disk
  // when the post is marshaled.
  UserId     string `json:"user_id,omitempty"`
  UserName   string `json:"user_name,omitempty"`
  UserAvatar string `json:"user_avatar,omitempty"`
  IsMyPost   bool   `json:"is_my_post,omitempty"`
  IsRepost   bool   `json:"is_repost,omitempty"`

  // Properties that are only used by the backend and thus don't have
  // to be marshaled to JSON for either the frontend or disk.
  FullPath string `json:"-"`
}

/**
 * InitNew()
 *
 * Initialize a new post creating the id and other fields.
 */
func (post *Post) InitNew(postData Post, userPath, userId, userName, userAvatar string, isMyUser bool) {
  // Create a new UUID
  uuidBytes, _ := uuid.NewV4()
  uuid := fmt.Sprintf("%s", uuidBytes)

  // Get the timestamp.
  created := time.Now().UnixNano()

  // The full path to the post.
  fullPath := path.Join(userPath, "posts", fmt.Sprintf("%d-post-%s.json", created, uuid))

  post.Id = uuid
  post.Title = postData.Title
  post.RepostPostId = postData.RepostPostId
  post.RepostUserName = postData.RepostUserName
  post.RepostUserAvatar = postData.RepostUserAvatar
  post.IsRepost = len(postData.RepostPostId) != 0
  post.Created = created
  post.UserId = userId
  post.UserName = userName
  post.UserAvatar = userAvatar
  post.IsMyPost = isMyUser
  post.FullPath = fullPath
}

/**
 * InitFromJson()
 *
 * Initialize a new post from json data from disk.
 */
func (post *Post) InitFromJson(rawJson []byte, fullPath string, userId string, userName string, userAvatar string, isMyUser bool) error {
  if err := json.Unmarshal(rawJson, post); err != nil {
    return err
  }
  post.UserId = userId
  post.UserName = userName
  post.UserAvatar = userAvatar
  post.IsMyPost = isMyUser
  post.FullPath = fullPath
  post.IsRepost = len(post.RepostPostId) != 0

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
  postClone := *post
  postClone.UserId = ""
  postClone.UserName = ""
  postClone.UserAvatar = ""
  postClone.IsMyPost = false
  postClone.IsRepost = false

  rawJson, err := json.Marshal(postClone)
  if err != nil {
    return err
  }

  return Write(postClone.FullPath, rawJson)
}

func (post *Post) Delete() error {
  return Delete(post.FullPath)
}

/**
 * Collection()
 *
 * Return a post collection wrapping this user.
 */
func (post *Post) Collection() *PostCollection {
  return &PostCollection{[]Post{*post}}
}

/**
 * Container()
 *
 * Return a post container wrapping this post.
 */
func (post *Post) Container() *PostContainer {
  return &PostContainer{*post}
}