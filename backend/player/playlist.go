package player

import (
	"context"
	"fmt"
	"mango/backend/catalog"
	"mango/backend/utils"
	"sync"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Playlist struct {
	ID      string
	Tracks  []*catalog.Track
	Current int
	Player  *Player
}

var (
	playlists     = make(map[string]*Playlist)
	playlistMutex sync.Mutex
)

func NewPlaylist(tracks []*catalog.Track) *Playlist {
	playlistMutex.Lock()
	defer playlistMutex.Unlock()
	id := fmt.Sprintf("%d", len(playlists)+1)
	pl := &Playlist{ID: id, Tracks: tracks}
	playlists[id] = pl
	return pl
}

func GetPlaylist(id string) (*Playlist, bool) {
	playlistMutex.Lock()
	defer playlistMutex.Unlock()
	pl, exists := playlists[id]
	return pl, exists
}

func (pl *Playlist) PlayCurrent(ctx context.Context) error {
	stopOtherPlaylists()
	track := pl.Tracks[pl.Current]
	streamer, format, err := decodeTrack(track.Filepath)
	if err != nil {
		return err
	}
	done := make(chan bool)
	pl.Player = newPlayer(streamer, format.SampleRate, done)
	pl.Player.play()
	return pl.handlePlayback(ctx, streamer, format, done)
}

func (pl *Playlist) handlePlayback(ctx context.Context, streamer beep.StreamSeekCloser, format beep.Format, done chan bool) error {
	ctrl := make(chan string)
	setupControlEvents(ctx, pl, ctrl)
	notifyTrackStart(ctx, pl)
	for {
		select {
		case <-done:
			return pl.NextTrack(ctx)
		case <-time.After(time.Second):
			updateCurrentPosition(ctx, streamer, format, pl)
		case request := <-ctrl:
			handlePlaylistControl(request, pl, ctx)
		}
	}
}

func setupControlEvents(ctx context.Context, pl *Playlist, ctrl chan string) {
	runtime.EventsOn(ctx, "ctrl:request", func(optionalData ...interface{}) {
		if len(optionalData) > 1 {
			handleCtrlRequest(optionalData, pl, ctrl)
		}
	})
}

func handleCtrlRequest(optionalData []interface{}, pl *Playlist, ctrl chan string) {
	request, ok := optionalData[0].(string)
	playlist, validPlaylist := optionalData[1].(string)

	if ok && validPlaylist && utils.IsValidCtrlRequest(request) && playlist == pl.ID {
		ctrl <- request
		ctrl <- playlist
	}
}

func notifyTrackStart(ctx context.Context, pl *Playlist) {
	runtime.EventsEmit(ctx, "track:playing", pl.Tracks[pl.Current], pl.Current)
	runtime.EventsEmit(ctx, "second:passed", 0, pl.ID)
}

func updateCurrentPosition(ctx context.Context, streamer beep.StreamSeekCloser, format beep.Format, pl *Playlist) {
	speaker.Lock()
	defer speaker.Unlock()
	runtime.EventsEmit(ctx, "second:passed", format.SampleRate.D(streamer.Position()).Round(time.Second), pl.ID)
}

func handlePlaylistControl(request string, pl *Playlist, ctx context.Context) {
	switch request {
	case "pause":
		playlists[pl.ID].Player.Pause()
	case "resume":
		playlists[pl.ID].Player.Resume()
	case "next":
		playlists[pl.ID].NextTrack(ctx)
	case "previous":
		playlists[pl.ID].PreviousTrack(ctx)
	}
}

func stopOtherPlaylists() {
	for _, p := range playlists {
		if p.Player != nil {
			p.Player.Pause()
		}
	}
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
