package types

type LatLong struct {
	Lat float64
	Lng float64
}

type Movie struct {
	ID          int
	TmdbID      int
	Title       string
	Coordinates LatLong
	Location    string
	PosterPath  string
	Overview    string
}
