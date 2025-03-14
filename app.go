package main

import (
	"context"
	"log"
	"mango/backend/catalog"
	"mango/backend/player"
	"mango/backend/utils"
	"os"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/flac"
	"github.com/gopxl/beep/v2/speaker"
)

type App struct {
	ctx    context.Context
	Player player.Player
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	sr := beep.SampleRate(41000)
	speaker.Init(sr, beep.SampleRate.N(sr, time.Second/10))
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

func (a *App) PlaySong(t catalog.Track) {
	f, err := os.Open(t.Filepath)
	if err != nil {
		log.Fatal(err)
	}
	streamer, format, err := flac.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()
	player, err := player.NewPlayer(streamer, format.SampleRate)
	if err != nil {
		log.Fatal(err)
	}
	player.PlayTrack(t)
}

func (a *App) PauseSong(player player.Player) {
	player.Pause()
}
