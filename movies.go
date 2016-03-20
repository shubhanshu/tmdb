package tmdb

import (
	"net/url"
	"errors"
	"strconv"
)

// Movie - Represents movie metadata
type Movie struct {
	ID             int    `json:"id"`
	Title          string `json:"original_title"`
	Overview       string `json:"overview"`
	PosterPath     string `json:"poster_path"`
	ReleaseDate    string `json:"release_date"`
}

// SearchResponse response of search query
type SearchResponse struct {
	Page         int
	Results      []Movie
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

type MovieSearchRequest struct {
	Name string
	Page int
	IncludeAdult bool
}

func NewMovieSearchRequest(name string) *MovieSearchRequest {
	return &MovieSearchRequest{name, 1, true}
}


var moviesSearchAPI = &apiConfig{
	host:            "http://api.themoviedb.org/3",
	path:            "/search/movie",
}


func (r *MovieSearchRequest) params() (q url.Values) {
	q = make(url.Values)
	q.Set("query", r.Name)
	q.Set("page", string(r.Page))
	q.Set("include_adult", strconv.FormatBool(r.IncludeAdult))
	return
}

func (c *Client) SearchMovies(r *MovieSearchRequest) (*Movie, error) {
	if (r.Name == "") {
		return nil, errors.New("Movie name is required to search")
	}
	resp := new(SearchResponse)
	err := c.getJSON(moviesSearchAPI, r, resp)
	if (err != nil) {
		return nil, err
	}

	if resp.TotalResults == 0 {
		return nil, errors.New("Result not found")
	}
	return &resp.Results[0], nil
}

