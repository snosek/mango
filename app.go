package main

import (
	"context"
	"log"
	"mango/backend/catalog"

	"github.com/wailsapp/wails/v2/pkg/runtime"
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
	dirPath, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		log.Print("Problem reading directory.")
		return "", err
	}
	return dirPath, nil
}

func (a *App) GetAlbums(fp string) ([]string, error) {
	return catalog.FetchDirectories(fp)
}
