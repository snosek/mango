package catalog

import (
	"errors"
	"mango/backend/utils"
)

type Catalog struct {
	Albums   map[string]*Album
	Filepath string
}

func NewCatalog(fp string) (Catalog, error) {
	if fp == "" {
		return Catalog{}, errors.New("empty music directory path")
	}
	catalog := Catalog{Filepath: fp}
	dirs, err := utils.FetchDirectories(fp)
	if err != nil {
		return catalog, err
	}
	catalog.Albums = make(map[string]*Album)
	for _, dir := range dirs {
		album, err := NewAlbum(dir)
		if err != nil {
			continue
		}
		catalog.Albums[album.ID] = &album
	}
	return catalog, nil
}
