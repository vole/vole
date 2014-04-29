package store

import (
	"encoding/json"
	"strings"
)

/**
 * UserCollection.
 */
type UserCollection struct {
	Users []User `json:"users"`
}

/**
 * Json()
 *
 * Return a string json representation of the user collection.
 */
func (userCollection *UserCollection) Json() (string, error) {
	out, err := json.Marshal(userCollection)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

/**
 * Filter()
 *
 * Filter users based on a query.
 */
func (userCollection *UserCollection) Filter(query string) {
	users := make([]User, 0)
	for _, user := range userCollection.Users {
		if strings.Contains(strings.ToLower(user.Name), strings.ToLower(query)) {
			users = append(users, user)
		}
	}
	userCollection.Users = users
}

/**
 * GetEmptyUserCollection()
 *
 * Return an empty collection.
 */
func GetEmptyUserCollection() *UserCollection {
	return &UserCollection{[]User{}}
}
