package store

import (
	"strings"
)

/**
 * User.
 */
type User struct {
	// Properties that should be saved to disk.
	Id     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Bio    string `json:"bio"`
	Online bool   `json:"online"`

	// Properties that are used by Vole backend and frontend, but not saved to disk.
	IsMyUser bool   `json:"is_my_user,omitempty"`
	Email    string `json:"email,omitempty"`

	// Properties that are only used by the backend and thus don't have
	// to be marshaled to JSON for either the frontend or disk.
	DirName      string `json:"-"`
	FullPath     string `json:"-"`
	UserJsonPath string `json:"-"`
}

/**
 * UserCollection.
 */
type UserCollection []User

func (collection UserCollection) Len() int {
	return len(collection)
}

/**
 * Filter()
 *
 * Filter users based on a query.
 */
func (collection *UserCollection) Filter(query string) *UserCollection {
	newCollection := make(UserCollection, 0)
	for _, user := range *collection {
		if strings.Contains(strings.ToLower(user.Name), strings.ToLower(query)) {
			newCollection = append(newCollection, user)
		}
	}
	return &newCollection
}

/**
 * FindById()
 *
 * Filter users based on a query.
 */
func (collection *UserCollection) FindById(id string) int {
	for index, user := range *collection {
		if user.Id == id {
			return index
		}
	}
	return -1
}
