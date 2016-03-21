package stores

import (
	"github.com/hailocab/go-geoindex"
	"github.com/shubhanshu/go-movie-locations/types"
	"reflect"
	"testing"
)

var (
	leicester = &geoindex.GeoPoint{"1", 51.511291, -0.128242}
	charring  = &geoindex.GeoPoint{"5", 51.508359, -0.124803}
	lewisham  = &geoindex.GeoPoint{"12", 51.46532, -0.0134}
)

func TestAddAndNearest(t *testing.T) {
	store := NewLocationStore()
	store.Add(types.LatLong{charring.Plat, charring.Plon}, 5)
	store.Add(types.LatLong{leicester.Plat, leicester.Plon}, 1)
	store.Add(types.LatLong{lewisham.Plat, lewisham.Plon}, 12)
	nearestPoints := store.Nearest(types.LatLong{charring.Plat, charring.Plon})
	expectedPoints := &[]int{5, 1}
	if !reflect.DeepEqual(nearestPoints, expectedPoints) {
		t.Errorf("Expected nearest points to be %v but were instead %v", expectedPoints, nearestPoints)
	}
}
