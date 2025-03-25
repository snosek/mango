package catalog

import (
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

func NewTrack(fp string) Track {
	t := Track{Filepath: fp}
	t.populateMetadata()
	t.ID = utils.HashTitle(t.Title)
	return t
}

func (t *Track) populateMetadata() {
	tags, err := taglib.ReadTags(t.Filepath)
	if err != nil {
		return
	}

	t.Title = utils.FirstOrEmpty(tags[taglib.Title])
	t.Artist = tags[taglib.Artist]
	t.TrackNumber = parseTrackNumber(tags[taglib.TrackNumber])

	props, err := taglib.ReadProperties(t.Filepath)
	if err == nil {
		t.Length = props.Length
		t.SampleRate = props.SampleRate
	}
}

func parseTrackNumber(nums []string) uint {
	if len(nums) == 0 {
		return 0
	}
	num, err := strconv.Atoi(nums[0])
	if err != nil {
		return 0
	}
	return uint(num)
}
