package catalog

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"mango/backend/files"
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
	ModTime  string
}

func NewAlbum(fp string) (*Album, error) {
	album := Album{Filepath: fp}
	if err := album.setTracks(); err != nil {
		return nil, fmt.Errorf("failed setting album tracks: %v", err)
	}
	if err := album.populateMetadata(); err != nil {
		return nil, fmt.Errorf("failed populating album metadata: %v", err)
	}
	return &album, nil
}

func (a *Album) setTracks() error {
	tracks, err := a.FetchTracks()
	if err != nil {
		return fmt.Errorf("failed fetching album tracks: %v", err)
	}
	a.Tracks = SortTracks(tracks)
	return nil
}

func (a *Album) populateMetadata() error {
	if len(a.Tracks) == 0 {
		return fmt.Errorf("no tracks in album %v", a.Title)
	}
	tags, err := files.ReadTags(a.Tracks[0].Filepath)
	if err != nil {
		return fmt.Errorf("failed reading album tags: %v", err)
	}
	a.Title = files.FirstOrEmpty(tags[taglib.Album])
	a.Artist = files.FirstOrFallback(tags[taglib.AlbumArtist], tags[taglib.Artist])
	a.Genre = tags[taglib.Genre]
	a.Length = a.calculateLength()
	a.Cover = a.encodeCover()
	a.ModTime = files.GetModificationTime(a.Filepath)
	a.ID = strings.ToLower(a.Filepath) + a.ModTime
	for _, t := range a.Tracks {
		t.AlbumID = a.ID
	}
	return nil
}

func (a *Album) encodeCover() string {
	cover, err := files.ReadAlbumCover(a.Filepath)
	if err != nil {
		return ""
	}
	encodedCover := resizeCover(cover)
	return encodedCover
}

func resizeCover(cover image.Image) string {
	m := resize.Resize(300, 300, cover, resize.NearestNeighbor)
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, m, nil); err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

func (a Album) calculateLength() time.Duration {
	var total time.Duration
	for _, t := range a.Tracks {
		total += t.Length
	}
	return total
}

func (a Album) FetchTracks() ([]*Track, error) {
	trackPaths, err := files.FetchAudioFiles(a.Filepath)
	if err != nil {
		return nil, err
	}
	tracks := []*Track{}
	for optionalTrackNum, fp := range trackPaths {
		t := NewTrack(fp, optionalTrackNum)
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
