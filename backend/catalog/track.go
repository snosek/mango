package catalog

import (
	"fmt"
	"mango/backend/utils"
	"strconv"
	"time"

	"go.senan.xyz/taglib"
)

type Track struct {
	ID          string
	Title       string
	Artist      []string
	TrackNumber uint
	Length      time.Duration
	SampleRate  uint
	AlbumID     string
	Filepath    string
}

func NewTrack(fp string, optionalTrackNum int) Track {
	t := Track{Filepath: fp}
	t.populateMetadata(optionalTrackNum)
	t.ID = utils.HashTitle(t.Title + fmt.Sprintf("%v", t.TrackNumber))
	return t
}

func (t *Track) populateMetadata(optionalTrackNum int) {
	tags, err := taglib.ReadTags(t.Filepath)
	if err != nil {
		return
	}

	t.Title = utils.FirstOrEmpty(tags[taglib.Title])
	t.Artist = tags[taglib.Artist]
	t.TrackNumber = parseTrackNumber(tags[taglib.TrackNumber], optionalTrackNum)

	props, err := taglib.ReadProperties(t.Filepath)
	if err == nil {
		t.Length = props.Length
		t.SampleRate = props.SampleRate
	}
}

func parseTrackNumber(nums []string, optionalTrackNum int) uint {
	if len(nums) == 0 {
		return uint(optionalTrackNum)
	}
	num, err := strconv.Atoi(nums[0])
	if err != nil {
		return uint(optionalTrackNum)
	}
	return uint(num)
}
