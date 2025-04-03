package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"mango/backend/catalog"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func NewDB() (*DB, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}
	mangoDir := filepath.Join(homeDir, ".mango")
	if err := os.MkdirAll(mangoDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}
	dbPath := filepath.Join(mangoDir, "catalog.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err := initSchema(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize database schema: %w", err)
	}
	return &DB{DB: db}, nil
}

func initSchema(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS config (
			musicDirPath TEXT
		);

		CREATE TABLE IF NOT EXISTS albums (
			id TEXT PRIMARY KEY,
			title TEXT,
			artist TEXT,
			genre TEXT,
			length INTEGER,
			cover TEXT,
			filepath TEXT
		);

		CREATE TABLE IF NOT EXISTS tracks (
			id TEXT PRIMARY KEY,
			title TEXT,
			artist TEXT,
			track_number INTEGER,
			length INTEGER,
			sample_rate INTEGER,
			album_id TEXT,
			filepath TEXT,
			FOREIGN KEY(album_id) REFERENCES albums(id)
		);
	`)
	if err != nil {
		return err
	}
	return err
}

func (db *DB) Close() error {
	return db.Close()
}

func (db *DB) SaveCatalog(catalog *catalog.Catalog) error {
	tx, err := db.Begin()
	defer tx.Rollback()
	if err != nil {
		return err
	}
	for _, album := range catalog.Albums {
		if err := db.saveAlbum(tx, album); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (db *DB) saveAlbum(tx *sql.Tx, album *catalog.Album) error {
	fmt.Println("saving album " + album.Title)
	artistJSON, err := json.Marshal(album.Artist)
	if err != nil {
		return fmt.Errorf("failed to marshal artist: %w", err)
	}
	genreJSON, err := json.Marshal(album.Genre)
	if err != nil {
		return fmt.Errorf("failed to marshal genre: %w", err)
	}
	_, err = tx.Exec(`
		INSERT OR REPLACE INTO albums
		(id, title, artist, genre, length, cover, filepath)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`,
		album.ID,
		album.Title,
		string(artistJSON),
		string(genreJSON),
		album.Length.Nanoseconds(),
		album.Cover,
		album.Filepath,
	)
	if err != nil {
		return fmt.Errorf("failed to insert album: %w", err)
	}
	for _, track := range album.Tracks {
		if err := db.saveTrack(tx, track); err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) saveTrack(tx *sql.Tx, track *catalog.Track) error {
	artistJSON, err := json.Marshal(track.Artist)
	if err != nil {
		return fmt.Errorf("failed to marshal track artist: %w", err)
	}
	_, err = tx.Exec(`
		INSERT OR REPLACE INTO tracks
		(id, title, artist, track_number, length, sample_rate, album_id, filepath)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		track.ID,
		track.Title,
		string(artistJSON),
		track.TrackNumber,
		track.Length.Nanoseconds(),
		track.SampleRate,
		track.AlbumID,
		track.Filepath,
	)
	if err != nil {
		return fmt.Errorf("failed to insert track: %w", err)
	}
	return nil
}

func (db *DB) LoadCatalog() (*catalog.Catalog, error) {
	result := catalog.Catalog{
		Albums: make(map[string]*catalog.Album),
	}
	albums, err := db.loadAlbums()
	if err != nil {
		return nil, err
	}
	for _, album := range albums {
		tracks, err := db.loadTracksForAlbum(album.ID)
		if err != nil {
			return nil, err
		}
		album.Tracks = tracks
		result.Albums[album.ID] = album
	}
	return &result, nil
}

func (db *DB) loadAlbums() ([]*catalog.Album, error) {
	rows, err := db.Query(`SELECT id, title, artist, genre, length, cover, filepath FROM albums`)
	if err != nil {
		return nil, fmt.Errorf("failed to query albums: %w", err)
	}
	defer rows.Close()
	var albums []*catalog.Album
	for rows.Next() {
		album := &catalog.Album{}
		var artistJSON, genreJSON string
		var lengthNano int64
		err := rows.Scan(
			&album.ID,
			&album.Title,
			&artistJSON,
			&genreJSON,
			&lengthNano,
			&album.Cover,
			&album.Filepath,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan album row: %w", err)
		}
		err = json.Unmarshal([]byte(artistJSON), &album.Artist)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal artist: %w", err)
		}
		err = json.Unmarshal([]byte(genreJSON), &album.Genre)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal genre: %w", err)
		}
		album.Length = time.Duration(lengthNano)
		albums = append(albums, album)
	}
	return albums, nil
}

func (db *DB) loadTracksForAlbum(albumID string) ([]*catalog.Track, error) {
	rows, err := db.Query(`
		SELECT id, title, artist, track_number, length, sample_rate, album_id, filepath
		FROM tracks
		WHERE album_id = ?
		ORDER BY track_number
	`, albumID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tracks: %w", err)
	}
	defer rows.Close()
	var tracks []*catalog.Track
	for rows.Next() {
		track := &catalog.Track{}
		var artistJSON string
		var lengthNano int64

		err := rows.Scan(
			&track.ID,
			&track.Title,
			&artistJSON,
			&track.TrackNumber,
			&lengthNano,
			&track.SampleRate,
			&track.AlbumID,
			&track.Filepath,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan track row: %w", err)
		}
		err = json.Unmarshal([]byte(artistJSON), &track.Artist)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal track artist: %w", err)
		}
		track.Length = time.Duration(lengthNano)
		tracks = append(tracks, track)
	}
	return tracks, nil
}
