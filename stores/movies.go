// In-memory index of all movies and locations
package stores

import "github.com/shubhanshu/movie-locations/types"

type MovieStore struct {
	cache map[int]types.Movie
}

func NewMovieStore() *MovieStore {
	return &MovieStore{make(map[int]types.Movie)}
}

func (store *MovieStore) Add(movie *types.Movie) {
	store.cache[movie.ID] = *movie
}

func (store *MovieStore) Get(id int) types.Movie {
	return store.cache[id]
}
