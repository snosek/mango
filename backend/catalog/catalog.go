package catalog

import (
	"mango/backend/utils"
)

type Catalog struct {
	Albums   map[string]*Album
	Filepath string
}

func NewCatalog(fp string) (Catalog, error) {
	catalog := Catalog{Filepath: fp}
	directories, err := utils.FetchDirectories(fp)
	if err != nil {
		return catalog, err
	}
	catalog.Albums = make(map[string]*Album)
	for _, dir := range directories {
		album, err := NewAlbum(dir)
		if err != nil {
			continue
		}
		catalog.Albums[album.ID] = &album
	}
	return catalog, nil
}
