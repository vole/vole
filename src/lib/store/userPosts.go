/**
 * Parts of the user struct that provide access to the user's posts.
 */

package store

import (
	"encoding/json"
	"errors"
	"path"
	"sort"
	"strings"
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
 * NewPostFromJson()
 *
 * Called by POST requests from the frontend.
 */
func (user *User) NewPostFromJson(rawJson []byte) (*Post, error) {
	var post Post
	if err := json.Unmarshal(rawJson, &post); err != nil {
		return nil, err
	}
	newPost := user.NewPost(post.Title)
	return newPost, nil
}

/**
 * NewPost(title)
 *
 * Return a new Post struct for this user.
 */
func (user *User) NewPost(title string) *Post {
	post := &Post{}
	post.InitNew(title, user.FullPath, user.Id, user.Name, user.Avatar, user.IsMyUser)
	return post
}

/**
 * GetPosts()
 *
 * Get all of a user's posts.
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
		if err := post.InitFromJson(data, fullPath, user.Id, user.Name, user.Avatar, user.IsMyUser); err != nil {
			return nil, errors.New("No post or invalid json.")
		}
		collection = append(collection, post)
	}
	postCol := &PostCollection{collection}
	sort.Sort(postCol)
	return postCol, nil
}

/**
 * GetDrafts()
 *
 * Get all of a user's drafts.
 */
func (user *User) GetDrafts() (*PostCollection, error) {
	collection := make([]Post, 0)
	postFiles, _ := ReadDir(user.FullPath, "drafts")

	for _, postFile := range postFiles {
		if !strings.HasSuffix(postFile.Name(), ".json") {
			continue
		}
		fullPath := path.Join(user.FullPath, "drafts", postFile.Name())
		data, err := ReadFile(fullPath)
		if err != nil {
			continue
		}
		post := Post{}
		if err := post.InitFromJson(data, fullPath, user.Id, user.Name, user.Avatar, user.IsMyUser); err != nil {
			return nil, errors.New("No post or invalid json.")
		}
		post.Draft = true
		collection = append(collection, post)
	}
	postCol := &PostCollection{collection}
	sort.Sort(postCol)
	return postCol, nil
}
