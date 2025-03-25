package player

import (
	"errors"
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

const sampleRate = beep.SampleRate(44100)
const bufferSize = time.Second / 7

func InitSpeaker() {
	sr := beep.SampleRate(sampleRate)
	speaker.Init(sr, sr.N(bufferSize))
}

type Player struct {
	streamer beep.Streamer
	ctrl     *beep.Ctrl
	volume   *effects.Volume
}

func NewPlayer(st beep.Streamer, sr beep.SampleRate, ch chan bool) *Player {
	resampled := resampleStreamer(st, sr, sampleRate)
	streamer := beep.Seq(resampled, beep.Callback(func() {
		ch <- true
	}))
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
		return beep.Resample(16, from, to, streamer)
	}
	return streamer
}

func decodeAudioFile(filePath, fileType string) (beep.StreamSeekCloser, beep.Format, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, beep.Format{}, err
	}
	var (
		streamer beep.StreamSeekCloser
		format   beep.Format
	)
	switch fileType {
	case "flac":
		streamer, format, err = flac.Decode(f)
	case "mp3":
		streamer, format, err = mp3.Decode(f)
	case "wav":
		streamer, format, err = wav.Decode(f)
	case "vorbis":
		streamer, format, err = vorbis.Decode(f)
	}
	if err != nil {
		return nil, beep.Format{}, err
	}
	return streamer, format, nil
}

func decodeTrack(filePath string) (beep.StreamSeekCloser, beep.Format, error) {
	switch filepath.Ext(filePath) {
	case ".flac":
		return decodeAudioFile(filePath, "flac")
	case ".mp3":
		return decodeAudioFile(filePath, "mp3")
	case ".wav":
		return decodeAudioFile(filePath, "wav")
	case ".ogg":
		return decodeAudioFile(filePath, "vorbis")
	default:
		return nil, beep.Format{}, errors.New("Unsupported file format.")
	}
}
