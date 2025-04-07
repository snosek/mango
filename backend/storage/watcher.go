package storage

import (
	"context"
	"mango/backend/utils"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	*fsnotify.Watcher
	db *DB
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
	go func() {
		SyncCatalogInRealTime(ctx, w, w.db)
	}()
}

func (w *Watcher) Close() error {
	return w.Close()
}
