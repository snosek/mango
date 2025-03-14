package main

import (
	"context"
	"mango/backend/catalog"
	"mango/backend/player"
	"mango/backend/utils"
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

func (a *App) Play(playlist *player.Playlist) {
	if err := playlist.PlayCurrent(); err != nil {
		return
	}
}

func (a *App) PauseSong(playlist player.Playlist) {
	playlist.Player.Pause()
}

func (a *App) ResumeSong(playlist player.Playlist) {
	playlist.Player.Resume()
}
