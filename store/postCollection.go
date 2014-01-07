package store

import (
  "encoding/json"
)

/**
 * PostCollection.
 */
type PostCollection struct {
  Posts []Post `json:"posts"`
}

/**
 * For sorting.
 */
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
 * Json()
 *
 * Return a string json representation of the posts.
 */
func (collection *PostCollection) Json() (string, error) {
  out, err := json.Marshal(collection)
  if err != nil {
    return "", err
  }
  return string(out), nil
}

/**
 * Find()
 *
 * Find a post within a collection and return its index.
 */
func (collection *PostCollection) Find(id string) int {
  for i, post := range collection.Posts {
    if post.Id == id {
      return i
    }
  }
  return -1
}

/**
 * Limit()
 *
 * Reduce the post collection to the specified limit.
 */
func (collection *PostCollection) Limit(limit int) {
  if limit > 0 && limit < len(collection.Posts) {
    collection.Posts = collection.Posts[0:limit]
  }
}

/**
 * BeforeId()
 *
 * Reduce the post collection to only posts before the specified ID.
 */
func (collection *PostCollection) BeforeId(id string) {
  if id == "" {
    return
  }
  i := collection.Find(id)
  if i == -1 {
    return
  }
  start := i + 1
  if start == len(collection.Posts) {
    collection.Posts = collection.Posts[i:i]
    return
  }
  collection.Posts = collection.Posts[start:]
}

/**
 * GetEmptyPostCollection()
 *
 * Return an empty collection.
 */
func GetEmptyPostCollection() *PostCollection {
  return &PostCollection{[]Post{}}
}
