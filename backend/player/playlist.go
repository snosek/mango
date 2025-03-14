package player

import (
	"fmt"
	"mango/backend/catalog"
	"os"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/flac"
)

type Playlist struct {
	Tracks  []*catalog.Track
	Current int
	Player  *Player
}

func NewPlaylist(tracks []*catalog.Track) *Playlist {
	return &Playlist{Tracks: tracks, Current: 0}
}

func (pl *Playlist) PlayCurrent() error {
	fmt.Println("------------------------------------------------------------------- entered playCurrent ")
	if pl.Player != nil {
		pl.Player.Pause()
	}
	streamer, format, err := decodeFLAC(pl.Tracks[pl.Current].Filepath)
	if err != nil {
		return err
	}
	fmt.Println("------------------------------------------------------------------- decoded flac ")
	resampled := resampleStreamer(streamer, format.SampleRate, beep.SampleRate(sampleRate))
	done := make(chan bool)
	pl.Player = NewPlayer(beep.Seq(resampled, beep.Callback(func() {
		done <- true
	})))
	pl.Player = NewPlayer(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	pl.Player.Play()
	go func() {
		<-done
		pl.NextTrack()
	}()
	return nil
}

func (pl *Playlist) NextTrack() error {
	if pl.Current < len(pl.Tracks)-1 {
		pl.Current++
		return pl.PlayCurrent()
	}
	return nil
}

func (pl *Playlist) PreviousTrack() error {
	if pl.Current > 0 {
		pl.Current--
		return pl.PlayCurrent()
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
