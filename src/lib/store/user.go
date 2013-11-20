package store

import (
  "encoding/json"
  "fmt"
  "github.com/vole/gouuid"
  "github.com/vole/gravatar"
  "os"
  "path"
)

/**
 * User.
 */
type User struct {
  // Properties that should be saved to disk.
  Id     string `json:"id"`
  Name   string `json:"name"`
  Avatar string `json:"avatar"`

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
 * InitNew()
 *
 * Initialize a new user creating the id and other fields.
 */
func (user *User) InitNew(name, email, storePath, version string) {
  // Create a new UUID
  uuidBytes, _ := uuid.NewV4()
  uuid := fmt.Sprintf("%s", uuidBytes)

  // Query Gravatar for the user's avatar URL.
  gravatarUrl := ""
  if email != "" {
    emailHash := gravatar.EmailHash(email)
    url := gravatar.GetAvatarURL("https", emailHash, gravatar.DefaultMysteryMan, 60)
    gravatarUrl = url.String()
  }

  user.Id = uuid
  user.Name = name
  user.Email = email
  user.Avatar = gravatarUrl
  user.IsMyUser = true

  // The location of the folder and object for this user, on disk.
  user.DirName = user.Name + "_" + user.Id
  user.FullPath = path.Join(storePath, "users", user.DirName, version)
  user.UserJsonPath = path.Join(user.FullPath, "user", "user.json")
}

/**
 * InitFromJson(json)
 *
 * Initialize this user from json loaded from disk.
 */
func (user *User) InitFromJson(rawJson []byte, subDir string, storePath string, version string) error {
  if err := json.Unmarshal(rawJson, user); err != nil {
    return err
  }
  user.DirName = subDir
  user.FullPath = path.Join(storePath, "users", user.DirName, version)
  user.UserJsonPath = path.Join(user.FullPath, "user", "user.json")

  return nil
}

func (user *User) FilePath() string {
  return path.Join(user.FullPath, "files")
}

/**
 * EnsureDirs()
 *
 * Creates appropriate folder structure for a user if it doesn't exist.
 * E.g. /Users/bobby/Vole/users/bobby-on-vole_xxx-xxx-xxx/v1
 */
func (user *User) EnsureDirs() error {
  // Create the user's user directory.
  if err := os.MkdirAll(path.Join(user.FullPath, "user"), 0755); err != nil {
    return err
  }

  // Create the user's posts directory.
  if err := os.MkdirAll(path.Join(user.FullPath, "posts"), 0755); err != nil {
    return err
  }

  // Create the user's files directory.
  if err := os.MkdirAll(path.Join(user.FullPath, "files"), 0755); err != nil {
    return err
  }

  return nil
}

/**
 * Save()
 *
 * Save user to disk.
 */
func (user *User) Save() error {
  if err := user.EnsureDirs(); err != nil {
    return err
  }

  // Before marshaling JSON for saving to disk, we set all properties
  // that should not be saved to empty, so they are ignored by marshaller.
  userClone := *user
  userClone.IsMyUser = false
  userClone.Email = ""

  rawJson, err := json.Marshal(userClone)
  if err != nil {
    return err
  }

  return Write(userClone.UserJsonPath, rawJson)
}

/**
 * Collection()
 *
 * Return a user collection wrapping this user.
 */
func (user *User) Collection() *UserCollection {
  return &UserCollection{[]User{*user}}
}

/**
 * Container()
 *
 * Return a user container wrapping this user.
 */
func (user *User) Container() *UserContainer {
  return &UserContainer{*user}
}
