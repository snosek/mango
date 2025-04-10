package storage

import (
	"context"
	"fmt"
	"mango/backend/catalog"
	"mango/backend/utils"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func SyncCatalog(db *DB, musicDirPath string) error {
	if musicDirPath == "" {
		return fmt.Errorf("empty music directory path")
	}
	existingCatalog, err := db.LoadCatalog()
	if err != nil {
		return fmt.Errorf("failed to load catalog: %w", err)
	}
	scannedAlbums, err := utils.FetchDirectories(musicDirPath)
	if err != nil {
		return err
	}
	scannedAlbumsIDPath := make(map[string]string)
	for _, albumPath := range scannedAlbums {
		albumModTime := utils.GetModificationTime(albumPath)
		scannedAlbumsIDPath[strings.ToLower(albumPath)+albumModTime] = albumPath
	}
	toAdd := make(map[string]string)
	toRemove := make(map[string]string)
	for id, path := range scannedAlbumsIDPath {
		if _, exists := existingCatalog.Albums[id]; !exists {
			toAdd[id] = path
		}
	}
	for id, album := range existingCatalog.Albums {
		if _, exists := scannedAlbumsIDPath[id]; !exists {
			toRemove[id] = album.Filepath
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	for id, _ := range toRemove {
		if err := db.RemoveAlbum(id); err != nil {
			return fmt.Errorf("Failed to remove album %s: %v", id, err)
		}
	}
	for _, fp := range toAdd {
		album, err := catalog.NewAlbum(fp)
		if err != nil {
			return err
		}
		if err := db.saveAlbum(tx, album); err != nil {
			return fmt.Errorf("failed to add album %s: %w", album.Title, err)
		}
	}
	return tx.Commit()
}

func SyncCatalogInRealTime(ctx context.Context, w *Watcher) {
	w.Watch()
	go func() {
		for event := range w.AlbumEvents {
			w.syncEvent(event, ctx)
		}
	}()
}

func (w *Watcher) syncEvent(event AlbumEvent, ctx context.Context) {
	tx, err := w.db.Begin()
	if err != nil {
		fmt.Printf("failed to start db transaction: %v\n", err)
		return
	}
	switch event.Type {
	case "add":
		album, err := catalog.NewAlbum(event.Path)
		if err != nil {
			fmt.Printf("failed to create new album: %v\n", err)
			return
		}
		if err := w.db.saveAlbum(tx, album); err != nil {
			tx.Rollback()
			fmt.Printf("failed to save album: %v\n", err)
			return
		}
		if err := tx.Commit(); err != nil {
			fmt.Printf("failed to commit db transaction: %v\n", err)
			return
		}
		runtime.EventsEmit(ctx, "album:addedOrRemoved")
	case "remove":
		if err := w.db.RemoveAlbumByPath(event.Path); err != nil {
			tx.Rollback()
			fmt.Printf("failed to remove album: %v\n", err)
			return
		}
		if err := tx.Commit(); err != nil {
			fmt.Printf("failed to commit db transaction: %v\n", err)
			return
		}
		runtime.EventsEmit(ctx, "album:addedOrRemoved")
	}
}
