package catalog

import (
	"mango/backend/utils"
	"sort"
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
	Tracks   []*Track
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
	album.Tracks = SortTracks(tracks)
	return album, nil
}

func (a Album) FetchTracks() ([]*Track, error) {
	trackPaths, err := utils.FetchAudioFiles(a.Metadata.Filepath)
	if err != nil {
		return nil, err
	}
	tracks := []*Track{}
	for _, fp := range trackPaths {
		t := NewTrack(fp)
		tracks = append(tracks, &t)
	}
	return tracks, nil
}

func SortTracks(tracksToSort []*Track) []*Track {
	sort.SliceStable(tracksToSort, func(i, j int) bool {
		return tracksToSort[i].TrackNumber < tracksToSort[j].TrackNumber
	})
	return tracksToSort
}
