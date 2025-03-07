package catalog

import (
	"errors"
	"log"
	"mango/backend/utils"
)

type AlbumMetadata struct {
	Filepath   string
	Title      string
	Artist     []string
	Genre      []string
	Length     int
	SampleRate int
}

type Album struct {
	Metadata *AlbumMetadata
	Tracks   map[string]*Track
}

func NewAlbum(fp string) Album {
	a := Album{
		Metadata: &AlbumMetadata{
			Filepath: fp,
		},
	}
	tracks, err := a.FetchTracks()
	log.Println(tracks)
	if err != nil {
		log.Println(err.Error())
	}
	a.Tracks = tracks
	return a
}

func (a Album) FetchTracks() (map[string]*Track, error) {
	trackPaths, err := utils.FetchAudioFiles(a.Metadata.Filepath)
	returnErr := errors.New("")
	if err != nil {
		returnErr = err
	}
	result := make(map[string]*Track, 0)
	for _, fp := range trackPaths {
		t := NewTrack(fp)
		result[t.Metadata.Title] = &t
	}
	log.Println(result)
	return result, returnErr
}
