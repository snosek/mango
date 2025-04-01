package player

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
	"github.com/gopxl/beep/v2/flac"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/vorbis"
	"github.com/gopxl/beep/v2/wav"
)

const (
	sampleRate = beep.SampleRate(44100)
	bufferSize = time.Second / 7
)

func InitSpeaker() {
	speaker.Init(sampleRate, sampleRate.N(bufferSize))
}

type Player struct {
	streamer beep.Streamer
	ctrl     *beep.Ctrl
	volume   *effects.Volume
}

func newPlayer(st beep.Streamer, sr beep.SampleRate, done chan bool) *Player {
	resampled := resampleStreamer(st, sr, sampleRate)
	streamer := beep.Seq(resampled, beep.Callback(func() { done <- true }))
	ctrl := &beep.Ctrl{Streamer: streamer, Paused: false}
	volume := &effects.Volume{Streamer: ctrl, Base: 2, Volume: 0}
	return &Player{
		streamer: streamer,
		ctrl:     ctrl,
		volume:   volume,
	}
}

func (p *Player) play() {
	speaker.Play(p.volume)
}

func (p *Player) Pause() {
	speaker.Lock()
	defer speaker.Unlock()
	p.ctrl.Paused = true
}

func (p *Player) Resume() {
	speaker.Lock()
	defer speaker.Unlock()
	p.ctrl.Paused = false
}

func (p *Player) setVolume(vol float64) {
	speaker.Lock()
	defer speaker.Unlock()
	p.volume.Volume = vol
}

func resampleStreamer(streamer beep.Streamer, from, to beep.SampleRate) beep.Streamer {
	if from != to {
		return beep.Resample(16, from, to, streamer)
	}
	return streamer
}

type decoderFunc func(io.Reader) (beep.StreamSeekCloser, beep.Format, error)

func decodeTrack(filePath string) (beep.StreamSeekCloser, beep.Format, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, beep.Format{}, err
	}
	decoders := map[string]decoderFunc{
		".flac": wrapDecoder(flac.Decode),
		".mp3":  wrapReadCloserDecoder(mp3.Decode),
		".wav":  wrapDecoder(wav.Decode),
		".ogg":  wrapReadCloserDecoder(vorbis.Decode),
	}
	decoder, ok := decoders[filepath.Ext(filePath)]
	if !ok {
		return nil, beep.Format{}, errors.New("unsupported file format")
	}
	return decoder(f)
}

func wrapDecoder(decode func(io.Reader) (beep.StreamSeekCloser, beep.Format, error)) decoderFunc {
	return func(r io.Reader) (beep.StreamSeekCloser, beep.Format, error) {
		return decode(r)
	}
}

func wrapReadCloserDecoder(decode func(io.ReadCloser) (beep.StreamSeekCloser, beep.Format, error)) decoderFunc {
	return func(r io.Reader) (beep.StreamSeekCloser, beep.Format, error) {
		readCloser, ok := r.(io.ReadCloser)
		if !ok {
			readCloser = io.NopCloser(r)
		}
		return decode(readCloser)
	}
}
