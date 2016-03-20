package tmdb

import (
	"net/http/httptest"
	"fmt"
	"net/http"
	"testing"
	"reflect"
)

const apiKey = "apikey"


func mockServer(status int, body string) *httptest.Server {
	server := new(httptest.Server)

	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		fmt.Fprintln(w, body)
	}))

	return server
}


func TestSearchingFightClub(t *testing.T) {
	response := `{
	  "page": 1,
	  "results": [
	    {
	      "adult": false,
	      "id": 550,
	      "original_language": "en",
	      "original_title": "Fight Club",
	      "overview": "A ticking-time-bomb insomniac and a slippery soap salesman channel primal male aggression into a shocking new form of therapy.",
	      "release_date": "1999-10-14",
	      "poster_path": "/811DjJTon9gD6hZ8nCjSitaIXFQ.jpg"
	    }
	  ],
	  "total_pages": 1,
	  "total_results": 1
	}`

	server := mockServer(200, response)
	defer server.Close()

	request := NewMovieSearchRequest("fight club")
	client, _ := NewClient(WithAPIKey(apiKey))
	client.baseURL = server.URL
	movie, err := client.SearchMovies(request)
	if err != nil {
		t.Errorf("Get returned non nil error: %v", err)
	}

	correctMovie := &Movie {
		ID: 550,
		Title: "Fight Club",
		Overview: "A ticking-time-bomb insomniac and a slippery soap salesman channel primal male aggression into a shocking new form of therapy.",
		PosterPath: "/811DjJTon9gD6hZ8nCjSitaIXFQ.jpg",
		ReleaseDate: "1999-10-14",
	}

	if !reflect.DeepEqual(movie, correctMovie) {
		t.Errorf("expected %+v, was %+v", correctMovie, movie)
	}
}


func TestNoResult(t *testing.T) {
	response := ""

	server := mockServer(404, response)
	defer server.Close()

	request := NewMovieSearchRequest("fight club")
	client, _ := NewClient(WithAPIKey(apiKey))
	client.baseURL = server.URL
	movie, err := client.SearchMovies(request)
	if err == nil {
		t.Errorf("Get did not return an error. Instead it returned: %v", movie)
	}
}

func TestMissingResult(t *testing.T) {
	response := `{
	  "page": 1,
	  "results": [
	  ],
	  "total_pages": 1,
	  "total_results": 0
	}`

	server := mockServer(200, response)
	defer server.Close()

	request := NewMovieSearchRequest("fight club")
	client, _ := NewClient(WithAPIKey(apiKey))
	client.baseURL = server.URL
	movie, err := client.SearchMovies(request)
	if err == nil {
		t.Errorf("Get did not return an error. Instead it returned: %v", movie)
	}
}

func TestMissingQuery(t *testing.T) {
	server := mockServer(200, "")
	defer server.Close()

	request := NewMovieSearchRequest("")
	client, _ := NewClient(WithAPIKey(apiKey))
	client.baseURL = server.URL
	movie, err := client.SearchMovies(request)
	if err == nil {
		t.Errorf("Get did not return an error. Instead it returned: %v", movie)
	}
}