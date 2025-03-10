package catalog

import (
	"mango/backend/utils"
	"sort"
	"time"

	"go.senan.xyz/taglib"
)

type Album struct {
	Title    string
	Artist   []string
	Genre    []string
	Length   time.Duration
	Tracks   []*Track
	Filepath string
}

func NewAlbum(fp string) (Album, error) {
	album := Album{
		Filepath: fp,
	}
	tracks, err := album.FetchTracks()
	if err != nil {
		return album, err
	}
	album.Tracks = SortTracks(tracks)
	album.SetMetadata()
	return album, nil
}

func (a *Album) SetMetadata() {
	tags, err := taglib.ReadTags(a.Tracks[0].Filepath)
	if err != nil {
		return
	}
	if tags[taglib.Album] != nil {
		a.Title = tags[taglib.Album][0]
	}
	a.Artist = tags[taglib.AlbumArtist]
	a.Genre = tags[taglib.Genre]
	a.Length = a.GetAlbumLength()
}

func (a Album) GetAlbumLength() time.Duration {
	var albumLength time.Duration
	for _, t := range a.Tracks {
		albumLength += t.Length
	}
	return albumLength
}

func (a Album) FetchTracks() ([]*Track, error) {
	trackPaths, err := utils.FetchAudioFiles(a.Filepath)
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
