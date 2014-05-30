package store

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Event struct {
	Type string
	Path string
}

func (event *Event) String() string {
	return fmt.Sprintf("%s: %s", event.Type, event.Path)
}

type Scanner struct {
	Dirs      []string
	Event     chan *Event
	lastMTime int64
	ticker    *time.Ticker
	files     map[string]bool
	scanning  bool
}

func NewScanner() *Scanner {
	scanner := new(Scanner)
	scanner.Event = make(chan *Event)
	scanner.lastMTime = time.Now().UnixNano()
	scanner.files = map[string]bool{}
	return scanner
}

func (scanner *Scanner) Add(dir string) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		scanner.files[path] = true
		return nil
	})

	scanner.Dirs = append(scanner.Dirs, dir)
}

func (scanner *Scanner) Start() {
	// Perform regular scans on a .25 second interval.
	scanner.ticker = time.NewTicker(time.Minute / 4)

	go func() {
		for _ = range scanner.ticker.C {
			scanner.Scan()
		}
	}()
}

// Scan the Root folder for file modifications.
// TODO(aaron): Better error handling.
func (scanner *Scanner) Scan() {
	if scanner.scanning {
		return
	}

	scanner.scanning = true
	now := time.Now()
	latest := map[string]bool{}

	for _, dir := range scanner.Dirs {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}

			latest[path] = true

			if _, exists := scanner.files[path]; !exists {
				scanner.Event <- &Event{"CREATED", path}
			} else if info.ModTime().UnixNano() > scanner.lastMTime {
				scanner.Event <- &Event{"MODIFIED", path}
			}

			return nil
		})
	}

	// Check for deleted files.
	for path, _ := range scanner.files {
		if _, exists := latest[path]; !exists {
			scanner.Event <- &Event{"DELETED", path}
		}
	}

	scanner.scanning = false
	scanner.lastMTime = now.UnixNano()
	scanner.files = latest
}

// TODO(aaron): Close the channel?
func (scanner *Scanner) Stop() {
	if scanner.ticker != nil {
		scanner.ticker.Stop()
	}

	close(scanner.Event)
}

//var scanner = Scanner{Root: "/tmp"}
//scanner.Start()
//ev := <-scanner.Event
