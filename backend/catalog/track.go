package catalog

import (
	"strconv"
	"time"

	"go.senan.xyz/taglib"
)

type TrackMetadata struct {
	Filepath    string
	Title       string
	Artist      []string
	Genre       []string
	TrackNumber uint
	Length      time.Duration
	SampleRate  uint
}

type Track struct {
	Metadata *TrackMetadata
}

func NewTrack(fp string) Track {
	t := Track{
		Metadata: &TrackMetadata{
			Filepath: fp,
		},
	}
	t.Metadata = t.FetchTrackMetadata()
	return t
}

func (t Track) FetchTrackMetadata() *TrackMetadata {
	meta := &TrackMetadata{Filepath: t.Metadata.Filepath}

	tags := t.FetchTrackTags()
	if tags != nil {
		meta.Title = tags["TITLE"][0]
		meta.Artist = tags["ARTIST"]
		meta.Genre = tags["GENRE"]
		if trackNum, err := strconv.Atoi(tags["TRACKNUMBER"][0]); err == nil {
			meta.TrackNumber = uint(trackNum)
		}
	}

	props := t.FetchTrackProperties()
	meta.Length = props.Length
	meta.SampleRate = props.SampleRate

	return meta
}

func (t Track) FetchTrackTags() map[string][]string {
	tags, err := taglib.ReadTags(t.Metadata.Filepath)
	if err != nil {
		return nil
	}
	return tags
}

func (t Track) FetchTrackProperties() taglib.Properties {
	props, err := taglib.ReadProperties(t.Metadata.Filepath)
	if err != nil {
		return taglib.Properties{}
	}
	return props
}
