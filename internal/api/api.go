// Package api that contains different structures and interface to connect to PokeAPI.
package api

const PokeAPIBaseURL = "https://pokeapi.co/api/v2"

type Endpoint interface {
	Path() string
	Get(url string) (any, error)
}

func EndpointsRegistry() map[string]Endpoint {
	return map[string]Endpoint{
		"areas":   NewLocationAreasAPI(),
		"area":    NewLocationAreaAPI(),
		"pokemon": NewPokemonAPI(),
	}
}
