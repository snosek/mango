package catalog

import (
	"bytes"
	"encoding/base64"
	"image/jpeg"
	"mango/backend/utils"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/KononK/resize"
	"go.senan.xyz/taglib"
)

type Album struct {
	ID       string
	Title    string
	Artist   []string
	Genre    []string
	Length   time.Duration
	Tracks   []*Track
	Cover    string
	Filepath string
}

func NewAlbum(fp string) (Album, error) {
	album := Album{Filepath: fp}
	tracks, err := album.FetchTracks()
	if err != nil {
		return album, err
	}
	album.Tracks = SortTracks(tracks)
	album.SetMetadata()
	album.ID = strings.ToLower(album.Title) + utils.Hash(album.Title)
	for _, t := range album.Tracks {
		t.AlbumID = album.ID
	}
	return album, nil
}

func (a *Album) SetMetadata() {
	if len(a.Tracks) == 0 {
		return
	}
	tags, err := taglib.ReadTags(a.Tracks[0].Filepath)
	if err != nil {
		return
	}
	if tags[taglib.Album] != nil {
		a.Title = tags[taglib.Album][0]
	}
	if tags[taglib.AlbumArtist] != nil {
		a.Artist = tags[taglib.AlbumArtist]
	} else {
		a.Artist = tags[taglib.Artist]
	}
	a.Genre = tags[taglib.Genre]
	a.Length = a.GetAlbumLength()
	a.Cover = a.EncodeCover()
}

func (a *Album) EncodeCover() string {
	file, err := os.Open(filepath.Join(a.Filepath, "folder.jpg"))
	if err != nil {
		return ""
	}
	defer file.Close()
	cover, err := jpeg.Decode(file)
	if err != nil {
		return ""
	}
	m := resize.Resize(300, 300, cover, resize.NearestNeighbor)
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, m, nil)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

func (a Album) GetAlbumLength() time.Duration {
	var totalLength time.Duration
	for _, t := range a.Tracks {
		totalLength += t.Length
	}
	return totalLength
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

func SortTracks(tracks []*Track) []*Track {
	sort.SliceStable(tracks, func(i, j int) bool {
		return tracks[i].TrackNumber < tracks[j].TrackNumber
	})
	return tracks
}
