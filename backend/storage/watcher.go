package storage

import (
	"log"
	"mango/backend/files"
	"os"

	"github.com/fsnotify/fsnotify"
)

type AlbumEvent struct {
	Path string
	Type string
}

type Watcher struct {
	*fsnotify.Watcher
	AlbumEvents chan AlbumEvent
	db          *DB
}

func NewWatcher(db *DB) (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	fp := GetMusicDirPath(db.DB)
	err = w.Add(fp)
	if err != nil {
		return nil, err
	}
	return &Watcher{Watcher: w, db: db}, nil
}

func (w *Watcher) Watch() {
	w.AlbumEvents = make(chan AlbumEvent)
	go w.watchLoop()
}

func (w *Watcher) watchLoop() {
	defer close(w.AlbumEvents)
	for {
		select {
		case event, ok := <-w.Events:
			if !ok {
				return
			}
			w.processEvent(event)
		case err, ok := <-w.Errors:
			if !ok {
				return
			}
			log.Printf("watcher error: %v", err)
		}
	}
}

func (w *Watcher) processEvent(event fsnotify.Event) {
	if files.IsSystemFile(event.Name) {
		return
	}
	switch {
	case event.Has(fsnotify.Write), event.Has(fsnotify.Create):
		info, err := os.Stat(event.Name)
		if err != nil {
			log.Printf("failed to stat created/written item: %v", err)
			return
		}
		if !info.IsDir() {
			return
		}
		w.AlbumEvents <- AlbumEvent{Path: event.Name, Type: "add"}
	case event.Has(fsnotify.Rename), event.Has(fsnotify.Remove):
		w.AlbumEvents <- AlbumEvent{Path: event.Name, Type: "remove"}
	}
}

func (w *Watcher) Close() error {
	return w.Watcher.Close()
}
