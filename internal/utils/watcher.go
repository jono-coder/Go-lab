package utils

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type EvenType int

const (
	Create EvenType = iota
	Write
	Remove
	Rename
	Chmod
)

type Event struct {
	Path string
	Type EvenType
}

type Config struct {
	Recursive    bool
	Debounce     time.Duration
	ErrorHandler func(error)
}

type Watcher struct {
	watcher *fsnotify.Watcher
	cfg     Config
	mu      sync.Mutex
	events  map[string]Event
}

func NewWatcher(cfg Config) (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &Watcher{
		watcher: w,
		cfg:     cfg,
		events:  make(map[string]Event),
	}, nil
}

func (w *Watcher) Add(path string) error {
	if !w.cfg.Recursive {
		return w.watcher.Add(path)
	}

	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return w.watcher.Add(path)
		}
		return nil
	})
}

func (w *Watcher) Watch(out chan<- Event) {
	debounce := w.cfg.Debounce
	if debounce == 0 {
		debounce = 50 * time.Millisecond
	}

	go func() {
		ticker := time.NewTicker(debounce)
		defer ticker.Stop()

		for {
			select {
			case fe, ok := <-w.watcher.Events:
				if !ok {
					return
				}
				w.processEvent(fe)

			case err, ok := <-w.watcher.Errors:
				if !ok {
					return
				}
				if w.cfg.ErrorHandler != nil {
					w.cfg.ErrorHandler(err)
				} else {
					log.Println("watcher error:", err)
				}

			case <-ticker.C:
				// flush debounced events
				w.mu.Lock()
				for _, ev := range w.events {
					out <- ev
				}
				w.events = make(map[string]Event)
				w.mu.Unlock()
			}
		}
	}()
}

func (w *Watcher) Close() error {
	return w.watcher.Close()
}

func (w *Watcher) processEvent(fe fsnotify.Event) {
	e := Event{}

	switch {
	case fe.Op&fsnotify.Create == fsnotify.Create:
		e.Type = Create
	case fe.Op&fsnotify.Write == fsnotify.Write:
		e.Type = Write
	case fe.Op&fsnotify.Remove == fsnotify.Remove:
		e.Type = Remove
	case fe.Op&fsnotify.Rename == fsnotify.Rename:
		e.Type = Rename
	case fe.Op&fsnotify.Chmod == fsnotify.Chmod:
		e.Type = Chmod
	}

	e.Path = fe.Name

	log.Println("event:", e)

	w.mu.Lock()
	w.events[e.Path] = e // store latest occurrence only
	w.mu.Unlock()
}
