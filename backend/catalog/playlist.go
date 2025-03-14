package catalog

type Playlist struct {
	Tracks []Track
}

func NewPlaylist(tracks []Track) Playlist {
	return Playlist{
		Tracks: tracks,
	}
}
