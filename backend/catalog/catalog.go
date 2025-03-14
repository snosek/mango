package catalog

import (
	"mango/backend/utils"
)

type Catalog struct {
	Albums   []*Album
	Filepath string
}

func NewCatalog(fp string) (Catalog, error) {
	catalog := Catalog{
		Filepath: fp,
	}
	directories, err := utils.FetchDirectories(fp)
	if err != nil {
		return catalog, err
	}
	for _, dir := range directories {
		album, err := NewAlbum(dir)
		if err != nil {
			return catalog, err
		}
		catalog.Albums = append(catalog.Albums, &album)
	}
	return catalog, nil
}
