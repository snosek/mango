package main

import (
	"context"
	"fmt"
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

func (a *App) Play(playlistID string) {
	if pl, exists := player.GetPlaylist(playlistID); exists {
		pl.PlayCurrent()
	}
}

func (a *App) PauseSong(playlistID string) {
	fmt.Println(playlistID)
	if pl, exists := player.GetPlaylist(playlistID); exists {
		pl.Player.Pause()
	}

}

func (a *App) ResumeSong(playlistID string) {
	fmt.Println(playlistID)
	if pl, exists := player.GetPlaylist(playlistID); exists {
		pl.Player.Resume()
	}
}

func (a *App) GetPlaylist(playlistID string) *player.Playlist {
	if pl, exists := player.GetPlaylist(playlistID); exists {
		return pl
	}
	return nil
}
