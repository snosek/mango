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
	a.Cover, err = a.encodeCover()
	if err != nil {
		return err
	}
	a.ModTime, err = files.GetModificationTime(a.Filepath)
	if err != nil {
		return fmt.Errorf("error getting modification time for %s: %w", a.Title, err)
	}
	a.ID = strings.ToLower(a.Filepath) + a.ModTime
	for _, t := range a.Tracks {
		t.AlbumID = a.ID
	}
	return nil
}

func (a *Album) encodeCover() (string, error) {
	cover, err := files.ReadAlbumCover(a.Filepath)
	if err != nil {
		return "", fmt.Errorf("error reading album cover for %s: %w", a.Title, err)
	}
	encodedCover, err := resizeCover(cover)
	if err != nil {
		return "", fmt.Errorf("error resizing album cover for %s: %w", a.Title, err)
	}
	return encodedCover, nil
}

func resizeCover(cover image.Image) (string, error) {
	m := resize.Resize(300, 300, cover, resize.NearestNeighbor)
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, m, nil); err != nil {
		return "", fmt.Errorf("error encoding album cover: %w", err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
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
