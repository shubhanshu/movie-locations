// Indexes list of locations in-memory and supports finding nearest points
package stores

import (
	"github.com/hailocab/go-geoindex"
	"github.com/shubhanshu/movie-locations/types"
	"log"
	"strconv"
)

type LocationStore struct {
	index *geoindex.PointsIndex
}

// NewLocationStore - initialize the location store
func NewLocationStore() *LocationStore {
	return &LocationStore{geoindex.NewPointsIndex(geoindex.Km(0.5))}
}

// Add - add a movie location to the location store
func (store *LocationStore) Add(point types.LatLong, id int) {
	location := geoindex.NewGeoPoint(strconv.Itoa(id), point.Lat, point.Lng)
	store.index.Add(location)
}

var (
	all = func(_ geoindex.Point) bool { return true }
)

// Nearest - find nearest movie shooting locations to the given point
func (store *LocationStore) Nearest(point types.LatLong) *[]int {
	locations := store.index.KNearest(geoindex.NewGeoPoint("1", point.Lat, point.Lng), 1000, geoindex.Km(1), all)
	result := make([]int, 0)
	for _, location := range locations {
		id, err := strconv.Atoi(location.Id())
		if err == nil {
			result = append(result, id)
		} else {
			log.Fatalf("Error converting location id to int %v", err)
		}

	}
	return &result
}
