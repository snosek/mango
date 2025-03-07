package catalog

import (
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

func NewAlbum(fp string) (Album, error) {
	album := Album{
		Metadata: &AlbumMetadata{
			Filepath: fp,
		},
	}
	tracks, err := album.FetchTracks()
	if err != nil {
		return album, err
	}
	album.Tracks = tracks
	return album, nil
}

func (a Album) FetchTracks() (map[string]*Track, error) {
	trackPaths, err := utils.FetchAudioFiles(a.Metadata.Filepath)
	if err != nil {
		return nil, err
	}
	tracks := make(map[string]*Track)
	for _, fp := range trackPaths {
		t := NewTrack(fp)
		tracks[t.Metadata.Title] = &t
	}
	return tracks, nil
}
