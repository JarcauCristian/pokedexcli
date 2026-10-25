package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Result struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationAreaResponse struct {
	Count   int      `json:"count"`
	Next    *string  `json:"next"`
	Prev    *string  `json:"previous"`
	Results []Result `json:"results"`
}

type LocationAreasAPI struct {
	NextURL *string
	PrevURL *string
}

func NewLocationAreasAPI() *LocationAreasAPI {
	return &LocationAreasAPI{}
}

func (l *LocationAreasAPI) Update(response LocationAreaResponse) {
	l.NextURL = response.Next
	l.PrevURL = response.Prev
}

func (l *LocationAreasAPI) Path() string {
	return PokeAPIBaseURL + "/" + "location-area"
}

func (l *LocationAreasAPI) Get(url string) (any, error) {
	res, err := http.Get(url)
	if err != nil {
		return LocationAreaResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return LocationAreaResponse{}, fmt.Errorf("expected 200 OK got %s", res.Status)
	}

	var response LocationAreaResponse
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&response); err != nil {
		return LocationAreaResponse{}, err
	}

	return response, nil
}
