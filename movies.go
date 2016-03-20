package tmdb

import (
	"errors"
	"net/url"
	"strconv"
)

// Movie - Represents movie metadata
type Movie struct {
	ID          int    `json:"id"`
	Title       string `json:"original_title"`
	Overview    string `json:"overview"`
	PosterPath  string `json:"poster_path"`
	ReleaseDate string `json:"release_date"`
}

// SearchResponse response of search query
type SearchResponse struct {
	Page         int
	Results      []Movie
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

// MovieSearchRequest request representing search query
type MovieSearchRequest struct {
	Query        string
	Page         int
	IncludeAdult bool
	Year         int
}

// NewMovieSearchRequest helper utility to create new search request
func NewMovieSearchRequest(name string) *MovieSearchRequest {
	return &MovieSearchRequest{name, 1, true, 0}
}

var moviesSearchAPI = &apiConfig{
	host: "http://api.themoviedb.org/3",
	path: "/search/movie",
}

func (r *MovieSearchRequest) params() (q url.Values) {
	q = make(url.Values)
	q.Set("query", r.Query)
	q.Set("page", string(r.Page))
	q.Set("include_adult", strconv.FormatBool(r.IncludeAdult))
	if r.Year != 0 {
		q.Set("year", string(r.Year))
	}
	return
}

// SearchMovies search for movies based on given request, returns an array of results
func (c *Client) SearchMovies(req *MovieSearchRequest) ([]Movie, error) {
	if req.Query == "" {
		return nil, errors.New("Movie name is required to search")
	}
	resp := new(SearchResponse)
	err := c.getJSON(moviesSearchAPI, req, resp)
	if err != nil {
		return nil, err
	}

	if resp.TotalResults == 0 {
		return nil, errors.New("Result not found")
	}
	return resp.Results[0:], nil
}
