package catalog

type Catalog struct {
	albums map[string]*Album
	tracks map[string]*Track
}
