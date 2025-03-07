package catalog

type AlbumMetadata struct {
	filepath   string
	title      string
	artist     []string
	genre      []string
	length     int
	sampleRate int
}

type Album struct {
	metadata *AlbumMetadata
	tracks   map[string]*Track
}
