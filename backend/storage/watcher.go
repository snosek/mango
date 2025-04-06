package storage

import (
	"mango/backend/utils"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	watcher  *fsnotify.Watcher
	filepath string
}

func NewWatcher() (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	db, err := NewDB()
	if err != nil {
		return nil, err
	}
	fp := utils.GetMusicDirPath(db.DB)
	err = w.Add(fp)
	if err != nil {
		return nil, err
	}
	return &Watcher{watcher: w, filepath: fp}, nil
}

func (w *Watcher) Watch() {
}

func (w *Watcher) Close() error {
	return w.watcher.Close()
}
