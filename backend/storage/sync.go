package storage

import (
	"context"
	"fmt"
	"log"
	"mango/backend/catalog"
	"mango/backend/utils"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
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
	for id := range toRemove {
		_, err = tx.Exec("DELETE FROM tracks WHERE album_id = ?", id)
		if err != nil {
			return fmt.Errorf("failed to delete tracks for album: %w", err)
		}
		_, err = tx.Exec("DELETE FROM albums WHERE id = ?", id)
		if err != nil {
			return fmt.Errorf("failed to delete album: %w", err)
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

func SyncCatalogInRealTime(ctx context.Context, w *Watcher, db *DB) {
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
					if album, err := catalog.NewAlbum(event.Name); err == nil {
						if tx, err := db.Begin(); err == nil {
							if err := db.saveAlbum(tx, album); err != nil {
								tx.Rollback()
								log.Printf("Failed to save album %s: %v", album.Title, err)
								continue
							}
							if err := tx.Commit(); err == nil {
								runtime.EventsEmit(ctx, "album:addedOrRemoved")
								log.Printf("Album added: %s", album.Title)
							} else {
								log.Printf("Failed to commit transaction: %v", err)
							}
						} else {
							log.Printf("Failed to begin transaction: %v", err)
						}
					} else {
						log.Printf("Error creating album from %s: %v", event.Name, err)
					}
				} else if err != nil {
					log.Printf("Error accessing file %s: %v", event.Name, err)
				}
			case event.Op&fsnotify.Rename != 0:
				if tx, err := db.Begin(); err == nil {
					if err := db.RemoveAlbumByPath(event.Name); err != nil {
						tx.Rollback()
						log.Printf("Failed to remove album at path %s: %v", event.Name, err)
						continue
					}
					if err := tx.Commit(); err == nil {
						runtime.EventsEmit(ctx, "album:addedOrRemoved")
						log.Printf("Album removed: %s", event.Name)
					} else {
						log.Printf("Failed to commit transaction: %v", err)
					}
				} else {
					log.Printf("Failed to begin transaction: %v", err)
				}
			}
		case err, ok := <-w.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}
