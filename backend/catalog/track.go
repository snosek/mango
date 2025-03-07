package catalog

import (
	"log"
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
	t.SetTrackMetadata()
	return t
}

func (t Track) SetTrackMetadata() {
	meta := t.FetchTrackTags()
	t.Metadata.Title = meta["TITLE"][0]
	t.Metadata.Artist = meta["ARTIST"]
	t.Metadata.Genre = meta["GENRE"]
	trackNumber, err := strconv.Atoi(meta["TRACKNUMBER"][0])
	if err != nil {
		log.Print("Error parsing track number.")
	}
	t.Metadata.TrackNumber = uint(trackNumber)

	props := t.FetchTrackProperties()
	t.Metadata.Length = props.Length
	t.Metadata.SampleRate = props.SampleRate
}

func (t Track) FetchTrackTags() map[string][]string {
	tags, err := taglib.ReadTags(t.Metadata.Filepath)
	if err != nil {
		log.Print("Error reading tags: ", err.Error())
		return nil
	}
	return tags
}

func (t Track) FetchTrackProperties() taglib.Properties {
	props, err := taglib.ReadProperties(t.Metadata.Filepath)
	if err != nil {
		log.Print("Error reading properties: ", err.Error())
		return taglib.Properties{}
	}
	return props
}
