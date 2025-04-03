package storage

import (
	"fmt"
	"mango/backend/catalog"

	_ "github.com/mattn/go-sqlite3"
)

func SyncCatalog(db *DB, musicDirPath string) error {
	if musicDirPath == "" {
		return fmt.Errorf("empty music directory path")
	}
	existingCatalog, err := db.LoadCatalog()
	if err != nil {
		return fmt.Errorf("failed to load catalog: %w", err)
	}
	scannedCatalog, err := catalog.NewCatalog(musicDirPath)
	if err != nil {
		return fmt.Errorf("failed to scan music directory: %w", err)
	}
	toAdd := make(map[string]*catalog.Album)
	toRemove := make(map[string]*catalog.Album)
	for id, album := range scannedCatalog.Albums {
		if _, exists := existingCatalog.Albums[id]; !exists {
			toAdd[id] = album
		}
	}
	for id, album := range existingCatalog.Albums {
		if _, exists := scannedCatalog.Albums[id]; !exists {
			toRemove[id] = album
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
	for _, album := range toAdd {
		if err := db.saveAlbum(tx, album); err != nil {
			return fmt.Errorf("failed to add album %s: %w", album.Title, err)
		}
	}
	return tx.Commit()
}
