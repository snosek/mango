package catalog

import (
	"log"

	"go.senan.xyz/taglib"
)

type TrackMetadata struct {
	filepath    string
	title       string
	artist      []string
	genre       []string
	trackNumber int
	length      int
	sampleRate  int
}

type Track struct {
	metadata *TrackMetadata
}

func GetTrackInfo(fp string) map[string][]string {
	tags, err := taglib.ReadTags(fp)
	if err != nil {
		log.Print("Error reading metadata: ", err.Error())
		return nil
	}
	return tags
}
