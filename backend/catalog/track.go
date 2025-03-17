package catalog

import (
	"strconv"
	"time"

	"go.senan.xyz/taglib"
)

type Track struct {
	Title       string
	Artist      []string
	TrackNumber uint
	Length      time.Duration
	SampleRate  uint
	Cover       *string
	Filepath    string
}

func NewTrack(fp string) Track {
	t := Track{Filepath: fp}
	t.SetMetadata()
	return t
}

func (t *Track) SetMetadata() {
	tags, err := taglib.ReadTags(t.Filepath)
	if err != nil {
		return
	}
	if tags[taglib.Title] != nil {
		t.Title = tags[taglib.Title][0]
	}
	if tags[taglib.TrackNumber] != nil {
		if trackNum, err := strconv.Atoi(tags[taglib.TrackNumber][0]); err == nil {
			t.TrackNumber = uint(trackNum)
		}
	}
	t.Artist = tags[taglib.Artist]

	if props, err := taglib.ReadProperties(t.Filepath); err == nil {
		t.Length = props.Length
		t.SampleRate = props.SampleRate
	}
}
