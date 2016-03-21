package bootstrap

import (
	"github.com/gregjones/httpcache"
	"github.com/shubhanshu/tmdb"
	"net/http"
	"strconv"
	"github.com/callum-ramage/jsonconfig"
)

// NewSearchClient initialize a new tmdb client
func NewSearchClient() (*tmdb.Client, error) {
	config, err := jsonconfig.LoadAbstract("./conf.json", "")

	if err != nil {
		return nil, err
	}
	cachedTransport := httpcache.NewMemoryCacheTransport()
	cachedClient := &http.Client{Transport: cachedTransport}
	client, err := tmdb.NewClient(tmdb.WithAPIKey(config["tmdbApiKey"].String()), tmdb.WithHTTPClient(cachedClient))
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Search for a movie in tmdb
func Search(title string, year string, client *tmdb.Client) ([]tmdb.Movie, error) {
	intYear, err := strconv.Atoi(year)
	searchRequest := &tmdb.MovieSearchRequest{
		Query:        title,
		Page:         1,
		IncludeAdult: false,
		Year:         intYear,
	}
	movies, err := client.SearchMovies(searchRequest)
	if err != nil {
		return nil, err
	}
	return movies, nil
}
