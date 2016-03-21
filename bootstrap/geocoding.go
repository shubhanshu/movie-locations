package bootstrap

import (
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
	"github.com/callum-ramage/jsonconfig"
)

// NewGeoClient initialize a new maps client
func NewGeoClient() (*maps.Client, error) {
	config, err := jsonconfig.LoadAbstract("./conf.json", "")

	if err != nil {
		return nil, err
	}

	client, err := maps.NewClient(maps.WithAPIKey(config["mapsApiKey"].String()))
	if err != nil {
		return nil, err
	}
	return client, nil
}

// GeoCode - Find lat,longitude of location
func GeoCode(location string, client *maps.Client) (*maps.LatLng, error) {
	r := &maps.TextSearchRequest{
		Query: location,
	}
	result, err := client.TextSearch(context.Background(), r)
	if err != nil {
		return nil, err
	}
	return &result.Results[0].Geometry.Location, nil
}
