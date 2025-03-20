package player

import (
	"context"
	"fmt"
	"mango/backend/catalog"
	"os"
	"sync"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/flac"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Playlist struct {
	ID      string
	Tracks  []*catalog.Track
	Current int
	Player  *Player
}

var (
	playlists = make(map[string]*Playlist)
	mu        sync.Mutex
)

func NewPlaylist(tracks []*catalog.Track) *Playlist {
	id := fmt.Sprintf("%d", len(playlists)+1)
	pl := &Playlist{ID: id, Tracks: tracks}
	mu.Lock()
	playlists[id] = pl
	mu.Unlock()
	return pl
}

func GetPlaylist(id string) (*Playlist, bool) {
	mu.Lock()
	defer mu.Unlock()
	pl, exists := playlists[id]
	return pl, exists
}

func (pl *Playlist) PlayCurrent(ctx context.Context) error {
	for _, p := range playlists {
		if p.Player != nil {
			p.Player.Pause()
		}
	}
	currentTrack := pl.Tracks[pl.Current]
	streamer, format, err := decodeFLAC(currentTrack.Filepath)
	if err != nil {
		return err
	}
	resampled := resampleStreamer(streamer, format.SampleRate, sampleRate)
	done := make(chan bool)
	pl.Player = NewPlayer(beep.Seq(resampled, beep.Callback(func() {
		done <- true
	})))
	runtime.EventsEmit(ctx, "track:playing", currentTrack, pl.Current)
	pl.Player.Play()
	go func() {
		<-done
		pl.NextTrack(ctx)
	}()
	return nil
}

func (pl *Playlist) NextTrack(ctx context.Context) error {
	if pl.Current < len(pl.Tracks)-1 {
		pl.Current++
		return pl.PlayCurrent(ctx)
	}
	return nil
}

func (pl *Playlist) PreviousTrack(ctx context.Context) error {
	if pl.Current > 0 {
		pl.Current--
		return pl.PlayCurrent(ctx)
	}
	return nil
}

func decodeFLAC(filePath string) (beep.StreamSeekCloser, beep.Format, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, beep.Format{}, err
	}
	streamer, format, err := flac.Decode(f)
	if err != nil {
		return nil, beep.Format{}, err
	}
	return streamer, format, nil
}
