package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type LocationAreaAPI struct {
	name string
}

func NewLocationAreaAPI() *LocationAreaAPI {
	return &LocationAreaAPI{}
}

func (l *LocationAreaAPI) Path() string {
	return PokeAPIBaseURL + "/location-area/" + l.name
}

func (l *LocationAreaAPI) Update(name string) {
	l.name = name
}

func (l *LocationAreaAPI) Get(url string) (any, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected 200 OK got %s", res.Status)
	}

	var response struct {
		PokemonEncounters []struct {
			Pokemon struct {
				Name string `json:"name"`
			} `json:"pokemon"`
		} `json:"pokemon_encounters"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decoding location area: %w", err)
	}

	names := make([]string, 0, len(response.PokemonEncounters))
	for _, e := range response.PokemonEncounters {
		names = append(names, e.Pokemon.Name)
	}
	return names, nil
}
