package player

import (
	"fmt"
	"mango/backend/catalog"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
	"github.com/gopxl/beep/v2/speaker"
)

type Player struct {
	sampleRate beep.SampleRate
	streamer   beep.StreamSeeker
	ctrl       *beep.Ctrl
	resampler  *beep.Resampler
	volume     *effects.Volume
}

func NewPlayer(streamer beep.StreamSeeker, sampleRate beep.SampleRate) (*Player, error) {
	loopStreamer, err := beep.Loop2(streamer)
	if err != nil {
		return nil, err
	}
	ctrl := &beep.Ctrl{Streamer: loopStreamer, Paused: false}
	resampler := beep.ResampleRatio(6, 1, ctrl)
	volume := &effects.Volume{Streamer: resampler, Base: 2}
	return &Player{
		sampleRate: sampleRate,
		streamer:   streamer,
		ctrl:       ctrl,
		resampler:  resampler,
		volume:     volume,
	}, nil
}

func (p Player) PlayTrack(t catalog.Track) {
	fmt.Print("playing...")
	done := make(chan bool)
	p.resampler.SetRatio(float64(t.SampleRate) / float64(41000))
	speaker.Play(beep.Seq(p.volume, beep.Callback(func() {
		done <- true
	})))
	<-done
}

func (p Player) Pause() {
	speaker.Lock()
	p.ctrl.Paused = !p.ctrl.Paused
	speaker.Unlock()
}
