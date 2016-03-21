package bootstrap

import (
	"github.com/shubhanshu/tmdb"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

const maxInt int = int(^uint(0) >> 1)

// Match given an array of movies, finds the one that best matches the title
func Match(movies []tmdb.Movie, title string) *tmdb.Movie {
	minimumDistance := maxInt
	bestMovie := movies[0]
	for _, movie := range movies {
		distance := levenshtein.DistanceForStrings(
			[]rune(movie.Title),
			[]rune(title),
			levenshtein.DefaultOptions)
		if distance < minimumDistance {
			bestMovie = movie
			minimumDistance = distance
		}
	}
	return &bestMovie
}
