package main

import (
	"context"
	"mango/backend/catalog"
	"mango/backend/utils"
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
	return utils.GetDirectories(fp)
}

func (a *App) GetTrackInfo(fp string) catalog.Track {
	t := catalog.NewTrack(fp)
	return t
}
