package catalog

import (
	"strconv"
	"time"

	"go.senan.xyz/taglib"
)

type Track struct {
	Filepath    string
	Title       string
	Artist      []string
	Genre       []string
	TrackNumber uint
	Length      time.Duration
	SampleRate  uint
}

func NewTrack(fp string) Track {
	t := Track{
		Filepath: fp,
	}
	t.FetchTrackMetadata()
	return t
}

func (t *Track) FetchTrackMetadata() {
	tags := t.FetchTrackTags()
	if tags == nil {
		return
	}
	if tags["TITLE"] != nil {
		t.Title = tags["TITLE"][0]
	}
	if tags["TRACKNUMBER"] != nil {
		if trackNum, err := strconv.Atoi(tags["TRACKNUMBER"][0]); err == nil {
			t.TrackNumber = uint(trackNum)
		}
	}
	t.Artist = tags["ARTIST"]
	t.Genre = tags["GENRE"]

	props := t.FetchTrackProperties()
	t.Length = props.Length
	t.SampleRate = props.SampleRate
}

func (t Track) FetchTrackTags() map[string][]string {
	tags, err := taglib.ReadTags(t.Filepath)
	if err != nil {
		return nil
	}
	return tags
}

func (t Track) FetchTrackProperties() taglib.Properties {
	props, err := taglib.ReadProperties(t.Filepath)
	if err != nil {
		return taglib.Properties{}
	}
	return props
}
