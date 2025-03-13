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
	Filepath    string
}

func NewTrack(fp string) Track {
	t := Track{
		Filepath: fp,
	}
	t.SetMetadata()
	return t
}

func (t *Track) SetMetadata() {
	tags := t.FetchTags()
	if tags == nil {
		return
	}
	if tags[taglib.Title] != nil {
		t.Title = tags[taglib.Title][0]
	}
	if tags[taglib.TrackNumber] != nil {
		trackNum, err := strconv.Atoi(tags[taglib.TrackNumber][0])
		if err == nil {
			t.TrackNumber = uint(trackNum)
		}
	}
	t.Artist = tags[taglib.Artist]

	props := t.FetchProperties()
	t.Length = props.Length
	t.SampleRate = props.SampleRate
}

func (t Track) FetchTags() map[string][]string {
	tags, err := taglib.ReadTags(t.Filepath)
	if err != nil {
		return nil
	}
	return tags
}

func (t Track) FetchProperties() taglib.Properties {
	props, err := taglib.ReadProperties(t.Filepath)
	if err != nil {
		return taglib.Properties{}
	}
	return props
}
