package main

import (
	"context"
	"mango/backend/catalog"
	"mango/backend/player"
	"mango/backend/utils"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	player.InitSpeaker()
	a.ctx = ctx
}

func (a *App) GetDirPath() (string, error) {
	return utils.GetDirPath(a.ctx)
}

func (a *App) GetAlbums(fp string) ([]string, error) {
	return utils.FetchDirectories(fp)
}

func (a *App) GetTrack(fp string) catalog.Track {
	t := catalog.NewTrack(fp)
	return t
}

func (a *App) GetAlbum(fp string) catalog.Album {
	album, err := catalog.NewAlbum(fp)
	if err != nil {
		return catalog.Album{}
	}
	return album
}

func (a *App) GetCatalog(fp string) catalog.Catalog {
	cat, err := catalog.NewCatalog(fp)
	if err != nil {
		return catalog.Catalog{}
	}
	return cat
}

func (a *App) NewPlaylist(tracks []*catalog.Track) *player.Playlist {
	return player.NewPlaylist(tracks)
}

func (a *App) Play(playlistID string) {
	pl, exists := player.GetPlaylist(playlistID)
	if !exists {
		return
	}
	err := pl.PlayCurrent()
	runtime.EventsEmit(a.ctx, "playerUpdated", pl.ID)
	if err != nil {
		return
	}
}

func (a *App) PauseSong(playlistID string) {
	pl, exists := player.GetPlaylist(playlistID)
	if !exists {
		return
	}
	pl.Player.Pause()
	runtime.EventsEmit(a.ctx, "playerUpdated", pl.ID)
}

func (a *App) ResumeSong(playlistID string) {
	pl, exists := player.GetPlaylist(playlistID)
	if !exists {
		return
	}
	pl.Player.Resume()
	runtime.EventsEmit(a.ctx, "playerUpdated", pl.ID)
}

func (a *App) GetPlaylist(playlistID string) *player.Playlist {
	pl, exists := player.GetPlaylist(playlistID)
	if !exists {
		return nil
	}
	return pl
}
