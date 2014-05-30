package store

import (
	"container/list"
	"encoding/json"
	"errors"
	"regexp"
)

// An index of the file system.
type BTSyncStoreIndex struct {
	Posts   *list.List
	Scanner *Scanner
}

// Start the indexer.
func (index *BTSyncStoreIndex) Start() {
	// Create a new file system scanner.
	index.Scanner = NewScanner()

	// Listen for scanner events.
	go func() {
		for {
			select {
			case ev := <-index.Scanner.Event:
				if ev.Type == "CREATED" || ev.Type == "MODIFIED" {
					index.AddPost(ev.Path)
				} else if ev.Type == "DELETED" {
					index.RemovePost(ev.Path)
				}
			}
		}
	}()

	// Start scanning the file system for changes.
	index.Scanner.Start()
}

// Adds a directory to the index's file system monitor.
func (index *BTSyncStoreIndex) Watch(dir string) {
	index.Scanner.Add(dir)
}

// Add a post to the index.
func (index *BTSyncStoreIndex) AddPost(file string) error {
	data, err := ReadFile(file)
	if err != nil {
		return err
	}

	post := Post{}
	if err := json.Unmarshal(data, &post); err != nil {
		return errors.New("No post or invalid json.")
	}

	// Insert the post into the index.
	for e := index.Posts.Front(); e != nil; e = e.Next() {
		if post.Created > e.Value.(Post).Created {
			index.Posts.InsertBefore(post, e)
			logger.Printf("Adding new post %s before %s", post.Id, e.Value.(Post).Id)
			break
		}
	}

	return nil
}

// Remove a post from the index.
func (index *BTSyncStoreIndex) RemovePost(file string) error {
	// HACK(aaron): Grab the post ID from the file name.
	re := regexp.MustCompile("[0-9]+-post-(.*).json")
	match := re.FindStringSubmatch(file)
	if match == nil {
		return errors.New("Error parsing Id from file name: " + file)
	}

	id := match[1]

	// Remove the post with the given id from the index.
	for e := index.Posts.Front(); e != nil; e = e.Next() {
		if id > e.Value.(Post).Id {
			index.Posts.Remove(e)
			logger.Printf("Removing post %s from index", id)
			break
		}
	}

	return nil
}
