package main

import (
	"context"
	"encoding/json"
	"mango/backend/catalog"
	"mango/backend/utils"
	"os"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
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
	album, _ := catalog.NewAlbum(fp)
	JSONAlbum, _ := json.Marshal(album)
	os.WriteFile("test/test.json", JSONAlbum, 0666)
	return album
}

func (a *App) GetCatalog(fp string) catalog.Catalog {
	cat, err := catalog.NewCatalog(fp)
	if err != nil {
		return catalog.Catalog{}
	}
	return cat
}
