package store

import (
  "encoding/json"
  "errors"
  "path"
  "sort"
  "strings"
)

/**
 * UserStore.
 */
type UserStore struct {
  Path    string
  Version string
}

/**
 * GetEmpty()
 *
 * Return an empty user struct.
 */
func (userStore *UserStore) GetEmptyUser() *User {
  return new(User)
}

/**
 * NewUserFromContainerJson()
 *
 * Called by POST requests from the frontend.
 */
func (userStore *UserStore) NewUserFromContainerJson(rawJson []byte) (*User, error) {
  var container UserContainer
  if err := json.Unmarshal(rawJson, &container); err != nil {
    return nil, err
  }
  user := userStore.NewUser(container.User.Name, container.User.Email)
  return user, nil
}

/**
 * NewUser(name, email)
 *
 * Return a new User struct.
 */
func (userStore *UserStore) NewUser(name, email string) *User {
  user := &User{}
  user.InitNew(name, email, userStore.Path, userStore.Version)
  return user
}

/**
 * GetMyUser()
 *
 * Get the user pointed at by the my_user file.
 */
func (userStore *UserStore) GetMyUser() (*User, error) {
  subDir, err := userStore.getMyUserSubDir()
  if err != nil {
    return nil, err
  }

  user, err := userStore.getUserFromDir(subDir)
  if err != nil {
    return nil, err
  }

  user.IsMyUser = true

  return user, nil
}

/**
 * GetUsers()
 *
 * Get all users.
 */
func (userStore *UserStore) GetUsers() (*UserCollection, error) {
  collection := make([]User, 0)

  myUserDir, err := userStore.getMyUserSubDir()
  if err != nil {
    return nil, err
  }

  userDirs, _ := ReadDir(userStore.Path, "users")

  for _, dir := range userDirs {
    user, err := userStore.getUserFromDir(dir.Name())
    if err != nil {
      continue
    }
    user.IsMyUser = (myUserDir == dir.Name())
    collection = append(collection, *user)
  }

  return &UserCollection{collection}, nil
}

/**
 * GetPosts()
 *
 * Get all posts from all users.
 */
func (userStore *UserStore) GetPosts(limit int) (*PostCollection, error) {
  users, err := userStore.GetUsers()
  if err != nil {
    return nil, err
  }

  collection := make([]Post, 0)
  for _, user := range users.Users {
    userPosts, err := user.GetPosts()
    if err != nil {
      continue
    }

    collection = append(collection, userPosts.Posts...)
  }

  postCol := &PostCollection{collection}
  sort.Sort(postCol)

  if limit >= len(postCol.Posts) {
    limit = len(postCol.Posts)
  }

  if limit > 0 {
    postCol.Posts = postCol.Posts[0:limit]
  }

  return postCol, nil
}

/**
 * GetPostsBeforeId(id, limit)
 *
 * Get posts occurring before the specified `id`, limited by `limit`.
 */
func (userStore *UserStore) GetPostsBeforeId(id string, limit int) (*PostCollection, error) {
  posts, err := userStore.GetPosts(0)
  if err != nil {
    return nil, err
  }

  i := posts.Find(id)

  if i == -1 {
    return nil, errors.New("Id not present.")
  }

  start := i + 1
  end := start + limit

  if end >= len(posts.Posts) {
    end = len(posts.Posts) - 1
  }

  posts.Posts = posts.Posts[start:end]
  return posts, nil
}

/**
 * GetPostsAfterId(id, limit)
 *
 * Get posts occurring after the specified `id`, limited by `limit`.
 */
func (userStore *UserStore) GetPostsAfterId(id string, limit int) (*PostCollection, error) {
  posts, err := userStore.GetPosts(0)
  if err != nil {
    return nil, err
  }

  i := posts.Find(id)
  if i == -1 {
    return nil, errors.New("Id not present.")
  }

  end := i - 1
  start := end - limit

  if start < 0 {
    start = 0
  }

  posts.Posts = posts.Posts[start:end]
  return posts, nil
}

/**
 * SetMyUser(user)
 *
 * Set the 'my_user' file to point to the given user.
 */
func (userStore *UserStore) SetMyUser(user *User) error {
  if err := Write(path.Join(userStore.Path, "my_user"), []byte(user.DirName)); err != nil {
    return err
  }
  return nil
}

/**
 * getMyUserSubDir()
 *
 * Return the contents of the my_user file which is the sub dir containing
 * the current user.
 */
func (userStore *UserStore) getMyUserSubDir() (string, error) {
  rawId, err := ReadFile(userStore.Path, "my_user")
  if err != nil {
    return "", err
  }
  subDir := strings.TrimSpace(string(rawId))
  return subDir, nil
}

/**
 * getUserFromDir(subDir)
 *
 * Get the user from the given subdir
 */
func (userStore *UserStore) getUserFromDir(subDir string) (*User, error) {
  data, err := ReadFile(userStore.Path, "users", subDir, userStore.Version, "user", "user.json")
  if err != nil {
    return nil, err
  }

  user := &User{}
  if err := user.InitFromJson(data, subDir, userStore.Path, userStore.Version); err != nil {
    return nil, errors.New("No user or invalid json.")
  }

  return user, nil
}
