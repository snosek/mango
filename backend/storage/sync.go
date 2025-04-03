package storage

import (
	"fmt"
	"mango/backend/catalog"
	"mango/backend/utils"
	"strings"

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
