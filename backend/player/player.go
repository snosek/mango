package player

import (
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
	"github.com/gopxl/beep/v2/speaker"
)

const sampleRate = 44100
const bufferSize = time.Second / 10

func InitSpeaker() {
	sr := beep.SampleRate(sampleRate)
	speaker.Init(sr, sr.N(bufferSize))
}

type Player struct {
	streamer beep.Streamer
	ctrl     *beep.Ctrl
	volume   *effects.Volume
}

func NewPlayer(streamer beep.Streamer) *Player {
	ctrl := &beep.Ctrl{Streamer: streamer, Paused: false}
	volume := &effects.Volume{Streamer: ctrl, Base: 2, Volume: 0}
	return &Player{
		streamer: streamer,
		ctrl:     ctrl,
		volume:   volume,
	}
}

func (p *Player) Play() {
	speaker.Play(p.volume)
}

func (p *Player) Pause() {
	speaker.Lock()
	p.ctrl.Paused = true
	speaker.Unlock()
}

func (p *Player) Resume() {
	speaker.Lock()
	p.ctrl.Paused = false
	speaker.Unlock()
}

func (p *Player) SetVolume(vol float64) {
	speaker.Lock()
	p.volume.Volume = vol
	speaker.Unlock()
}

func resampleStreamer(streamer beep.Streamer, from, to beep.SampleRate) beep.Streamer {
	if from != to {
		return beep.Resample(10, from, to, streamer)
	}
	return streamer
}
