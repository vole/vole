package store

/**
 * Post.
 */
type Post struct {
	// Properties that should be saved to disk.
	Id       string `json:"id"`
	Body     string `json:"title"`
	Created  int64  `json:"created"`
	Modified int64  `json:"modified"`

	// Properties that are used by Vole backend and frontend, but not saved to disk
	// when the post is marshaled.
	User  *User `json:"user,omitempty"`
	Draft bool  `json:"draft,omitempty"`
}

/**
 * PostCollection.
 */
type PostCollection []Post

/**
 * For sorting.
 */
func (collection PostCollection) Len() int {
	return len(collection)
}

func (collection PostCollection) Less(i, j int) bool {
	return collection[i].Created > collection[j].Created
}

func (collection PostCollection) Swap(i, j int) {
	collection[i], collection[j] = collection[j], collection[i]
}

/**
 * FindById()
 *
 * Find a post within a collection and return its index.
 */
func (collection *PostCollection) FindById(id string) int {
	for i, post := range *collection {
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
func (collection PostCollection) Limit(limit int) {
	if limit > 0 && limit < collection.Len() {
		collection = collection[0:limit]
	}
}

/**
 * BeforeId()
 *
 * Reduce the post collection to only posts before the specified ID.
 */
func (collection PostCollection) BeforeId(id string) {
	if id == "" {
		return
	}
	i := collection.FindById(id)
	if i == -1 {
		return
	}
	start := i + 1
	if start == collection.Len() {
		collection = collection[i:i]
		return
	}
	collection = collection[start:]
}
