// Builds bootstrap data by first getting list of movies and shooting locations from SODA APIs
// It then decorates this with data from TMDB and Google Places API
// Finally it stores it in a json file
package main

import (
	"encoding/json"
	"github.com/shubhanshu/go-movie-locations/bootstrap"
	"github.com/shubhanshu/go-movie-locations/types"
	"github.com/shubhanshu/tmdb"
	"googlemaps.github.io/maps"
	"log"
	"os"
)

func main() {
	locations, err := bootstrap.Locations()
	if err != nil {
		log.Fatalf("Error fetching locations %v", err)
		return
	}
	id := 1
	moviedata := make([]types.Movie, 0)
	searchClient, err := bootstrap.NewSearchClient()
	if err != nil {
		log.Fatalf("Error fetching tmdb client %v", err)
		return
	}

	geoClient, err := bootstrap.NewGeoClient()
	if err != nil {
		log.Fatalf("Error fetching maps client %v", err)
		return
	}
	for i, l := range locations {
		if l.Location != "" {
			searcc, searche := search(l.Title, l.Year, searchClient)
			geoc, geoe := geoencode(l.Location, geoClient)
			moviedatum := types.Movie{}
			select {
			case movie := <-searcc:
				moviedatum.Title = movie.Title
				moviedatum.TmdbID = movie.ID
				moviedatum.Overview = movie.Overview
				moviedatum.PosterPath = movie.PosterPath
			case <-searche:
				continue
			}

			select {
			case latlng := <-geoc:
				moviedatum.Coordinates = types.LatLong{latlng.Lat, latlng.Lng}
				moviedatum.Location = l.Location
			case <-geoe:
				continue
			}
			moviedatum.ID = id
			moviedata = append(moviedata, moviedatum)
			id += 1
			if i%10 == 0 {
				log.Printf("Processed %d out of %d locations", id, len(locations))
			}

		}

	}
	storeResults(moviedata)
}

func storeResults(movieData []types.Movie) {
	os.MkdirAll("./data", 0777)
	file, err := os.OpenFile("data/movie_data.json", os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		log.Printf("Error creating file to store results %v. Try deleting the file", err)
	}
	defer file.Close()

	movie_json, err := json.MarshalIndent(movieData, "", "  ")
	if err != nil {
		log.Printf("Error marshalling data: %v", err.Error())
	}
	_, err = file.Write(movie_json)
	if err != nil {
		log.Printf("Error writing to file:%v", err.Error())
	}
}

func search(title string, year string, client *tmdb.Client) (<-chan *tmdb.Movie, <-chan error) {
	out := make(chan *tmdb.Movie)
	errc := make(chan error, 1)
	go func() {
		movieCandidates, err := bootstrap.Search(title, year, client)
		if err != nil {
			log.Printf("Error fetching movie metadata %v for movie %s", err, title)
			errc <- err
			return
		}
		out <- bootstrap.Match(movieCandidates, title)
		close(out)
	}()

	return out, errc
}

func geoencode(location string, client *maps.Client) (<-chan *maps.LatLng, <-chan error) {
	out := make(chan *maps.LatLng)
	errc := make(chan error, 1)
	go func() {
		latlng, err := bootstrap.GeoCode(location+", San Francisco", client)
		if err != nil {
			log.Printf("Error reverse geo encoding movie location %v for location %s", err, location)
			errc <- err
			return
		}
		out <- latlng
		close(out)
	}()
	return out, errc
}
