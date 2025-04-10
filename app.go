package main

import (
	"context"
	"mango/backend/catalog"
	"mango/backend/files"
	"mango/backend/player"
	"mango/backend/storage"
	"mango/backend/utils"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx       context.Context
	cat       catalog.Catalog
	DB        *storage.DB
	FSWatcher *storage.Watcher
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	player.InitSpeaker()
	a.ctx = ctx
	runtime.LogSetLogLevel(ctx, 3)
	err := a.NewDB()
	if err != nil {
		return
	}
	w, err := storage.NewWatcher(a.DB)
	if err != nil {
		return
	}
	a.FSWatcher = w
	storage.SyncCatalogInRealTime(a.ctx, a.FSWatcher)
}

func (a *App) NewDB() error {
	db, err := storage.NewDB()
	if err != nil {
		return err
	}
	a.DB = db
	return nil
}

func (a *App) SyncDB() {
	var musicDirPath string
	a.DB.QueryRow(`SELECT musicDirPath FROM config`).Scan(&musicDirPath)
	storage.SyncCatalog(a.DB, musicDirPath)
}

func (a *App) shutdown(ctx context.Context) {
	a.DB.Close()
	a.FSWatcher.Close()
}

func (a *App) GetDirPath() (string, error) {
	dirPath, err := utils.GetDirPath(a.ctx)
	if err != nil {
		return "", err
	}
	a.DB.Exec(`INSERT OR REPLACE INTO config (musicDirPath) VALUES (?)`, dirPath)
	return dirPath, nil
}

func (a *App) GetAlbums(fp string) ([]string, error) {
	return files.FetchDirectories(fp)
}

func (a *App) GetAlbum(fp string) catalog.Album {
	album, err := catalog.NewAlbum(fp)
	if err != nil {
		return catalog.Album{}
	}
	return *album
}

func (a *App) GetCatalog(fp string) catalog.Catalog {
	cat, err := catalog.NewCatalog(fp)
	if err != nil {
		return catalog.Catalog{}
	}
	a.cat = cat
	return cat
}

func (a *App) LoadCatalogFromDB() catalog.Catalog {
	cat, err := a.DB.LoadCatalog()
	if err != nil {
		return catalog.Catalog{}
	}
	return *cat
}

func (a *App) NewPlaylist(tracks []*catalog.Track) *player.Playlist {
	return player.NewPlaylist(tracks)
}

func (a *App) Play(playlistID string) {
	pl, exists := player.GetPlaylist(playlistID)
	if !exists {
		return
	}
	pl.PlayCurrent(a.ctx)
}

func (a *App) PauseSong(playlistID string) {
	if pl, exists := player.GetPlaylist(playlistID); exists {
		pl.Player.Pause()
	}
}

func (a *App) ResumeSong(playlistID string) {
	if pl, exists := player.GetPlaylist(playlistID); exists {
		pl.Player.Resume()
	}
}

func (a *App) PreviousTrack(playlistID string) {
	if pl, exists := player.GetPlaylist(playlistID); exists {
		pl.PreviousTrack(a.ctx)
	}
}

func (a *App) NextTrack(playlistID string) {
	if pl, exists := player.GetPlaylist(playlistID); exists {
		pl.NextTrack(a.ctx)
	}
}

func (a *App) GetPlaylist(playlistID string) *player.Playlist {
	if pl, exists := player.GetPlaylist(playlistID); exists {
		return pl
	}
	return nil
}

func (a *App) GetCurrentTrack(playlistID string) *catalog.Track {
	if pl, exists := player.GetPlaylist(playlistID); exists && pl.Current < len(pl.Tracks) {
		return pl.Tracks[pl.Current]
	}
	return nil
}

func (a *App) SaveCatalog(cat *catalog.Catalog) {
	a.DB.SaveCatalog(cat)
}

func (a *App) GetMusicDirPath() string {
	if a.DB == nil {
		return ""
	}
	return storage.GetMusicDirPath(a.DB.DB)
}

func (a *App) IsDBAvailable() bool {
	for {
		if a.DB != nil {
			return true
		}
	}
}
