package storage

import (
	"context"
	"log"
	"mango/backend/utils"
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
	fp := utils.GetMusicDirPath(db.DB)
	err = w.Add(fp)
	if err != nil {
		return nil, err
	}
	return &Watcher{Watcher: w, db: db}, nil
}

func (w *Watcher) Watch(ctx context.Context) {
	w.AlbumEvents = make(chan AlbumEvent)
	go func() {
		defer close(w.AlbumEvents)
		for {
			select {
			case event, ok := <-w.Events:
				if !ok {
					return
				}
				if utils.IsSystemFile(event.Name) {
					continue
				}
				switch {
				case event.Op&(fsnotify.Create|fsnotify.Write) != 0:
					if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
						w.AlbumEvents <- AlbumEvent{Path: event.Name, Type: "add"}
					}
				case event.Op&fsnotify.Rename != 0:
					w.AlbumEvents <- AlbumEvent{Path: event.Name, Type: "remove"}
				}
			case err, ok := <-w.Errors:
				if !ok {
					return
				}
				log.Printf("Watcher error: %v", err)
			}
		}
	}()
}

func (w *Watcher) Close() error {
	return w.Watcher.Close()
}
