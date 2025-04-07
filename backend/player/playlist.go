package player

import (
	"context"
	"fmt"
	"mango/backend/catalog"
	"mango/backend/utils"
	"math"
	"strconv"
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
	notifyTrackStart(ctx, pl)
	ctrl := setupControlEvents(ctx, pl)
	trackSwitch := setupTrackEvents(ctx, pl)
	for {
		select {
		case <-done:
			return pl.NextTrack(ctx)
		case <-time.After(time.Second):
			updateCurrentPosition(ctx, streamer, format, pl)
		case request := <-ctrl:
			handlePlaylistControl(ctx, request, pl, streamer, format, ctrl)
		case playlistID := <-trackSwitch:
			handleTrackSwitch(ctx, playlistID, pl, trackSwitch)
		}
	}
}

func setupControlEvents(ctx context.Context, pl *Playlist) chan string {
	ctrl := make(chan string)
	runtime.EventsOn(ctx, "ctrl:request", func(optionalData ...any) {
		if len(optionalData) > 1 {
			handleCtrlRequest(ctx, optionalData, pl, ctrl)
		}
	})
	return ctrl
}

func setupTrackEvents(ctx context.Context, pl *Playlist) chan string {
	trackSwitch := make(chan string)
	runtime.EventsOn(ctx, "track:switch", func(optionalData ...any) {
		if len(optionalData) > 1 {
			handleTrackSwitchRequest(ctx, optionalData, pl, trackSwitch)
		}
	})
	return trackSwitch
}

func handleCtrlRequest(ctx context.Context, optionalData []any, pl *Playlist, ctrl chan string) {
	request, validRequest := optionalData[0].(string)
	playlistID, validPlaylistID := optionalData[1].(string)
	if validRequest && validPlaylistID && utils.IsValidCtrlRequest(request) {
		ctrl <- request
		ctrl <- playlistID
	}
	if len(optionalData) > 2 {
		data, ok := optionalData[2].(string)
		if ok {
			ctrl <- data
		}
	}
}

func handleTrackSwitchRequest(ctx context.Context, optionalData []any, pl *Playlist, trackSwitch chan string) {
	playlistID, validPlaylistID := optionalData[0].(string)
	currentAlbum, validAlbum := optionalData[1].(string)
	trackNumber, validTrackNumber := optionalData[2].(string)
	if validAlbum && validPlaylistID && validTrackNumber {
		trackSwitch <- playlistID
		trackSwitch <- currentAlbum
		trackSwitch <- trackNumber
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

func handlePlaylistControl(ctx context.Context, request string, pl *Playlist, s beep.StreamSeekCloser, f beep.Format, ctrl chan string) {
	playlistID := <-ctrl
	if playlistID != pl.ID {
		return
	}
	switch request {
	case "pause":
		playlists[pl.ID].Player.Pause()
	case "resume":
		playlists[pl.ID].Player.Resume()
	case "next":
		playlists[pl.ID].NextTrack(ctx)
	case "previous":
		playlists[pl.ID].PreviousTrack(ctx)
	case "changePosition":
		position := <-ctrl
		positionFloat, err := strconv.ParseFloat(position, 64)
		if err != nil {
			return
		}
		positionFloat = math.Max(0, math.Min(1, positionFloat))
		currentTrack := pl.Tracks[pl.Current]
		totalSamples := f.SampleRate.N(currentTrack.Length)
		samplePosition := int(float64(totalSamples) * positionFloat)
		speaker.Lock()
		err = s.Seek(samplePosition)
		speaker.Unlock()
		updateCurrentPosition(ctx, s, f, pl)
	default:
		runtime.LogInfo(ctx, "Unknown request")
	}
}

func handleTrackSwitch(ctx context.Context, playlistID string, pl *Playlist, trackSwitch chan string) {
	if playlistID == "" || playlistID != pl.ID {
		return
	}
	albumID := <-trackSwitch
	trackNumber, err := strconv.Atoi(<-trackSwitch)
	if err != nil {
		return
	}
	if albumID == pl.Tracks[0].AlbumID {
		pl.Current = trackNumber
		pl.PlayCurrent(ctx)
	}
	if albumID != pl.Tracks[0].AlbumID {
		// tutaj tworzenie nowej playlisty i odtwarzanie jej. potrzebna lista trackow
		// wiec chyba najpierw trzeba zrobic storage
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
