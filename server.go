package main

import (
	"encoding/json"
	"github.com/shubhanshu/go-movie-locations/stores"
	"github.com/shubhanshu/go-movie-locations/types"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

type LocalStores struct {
	Locations *stores.LocationStore
	Movies    *stores.MovieStore
}

var (
	localStores LocalStores
)

func searchByLocation(w http.ResponseWriter, r *http.Request) {
	query := path.Base(r.URL.Path)
	latlong := strings.Split(query, ",")
	if len(latlong) != 2 {
		w.WriteHeader(http.StatusBadRequest)
	}
	lat, err := strconv.ParseFloat(latlong[0], 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	lng, err := strconv.ParseFloat(latlong[1], 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	ids := localStores.Locations.Nearest(types.LatLong{lat, lng})
	if len(*ids) == 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	movies := make([]types.Movie, 0)
	for _, id := range *ids {
		movies = append(movies, localStores.Movies.Get(id))
	}
	moviesj, err := json.MarshalIndent(movies, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(moviesj)
}

func main() {
	http.HandleFunc("/search/location/", searchByLocation)
	http.ListenAndServe(":8000", nil)
}

func init() {
	bootstrapFile, err := os.Open("data/movie_data.json")
	if err != nil {
		log.Fatalf("Error reading bootstrap data %v", err)
	}
	movies := make([]types.Movie, 0)
	decoder := json.NewDecoder(bootstrapFile)
	err = decoder.Decode(&movies)
	if err != nil {
		log.Fatalf("Error reading bootstrap data %v", err)
	}
	localStores = LocalStores{stores.NewLocationStore(), stores.NewMovieStore()}
	for _, movie := range movies {
		localStores.Movies.Add(&movie)
		localStores.Locations.Add(movie.Coordinates, movie.ID)
	}
	log.Println("Initialization done!")
}
