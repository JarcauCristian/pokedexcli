package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Stat struct {
	Name string `json:"name"`
}

type Statistic struct {
	Stat     Stat `json:"stat"`
	BaseStat int  `json:"base_stat"`
}

type Type struct {
	Type Result `json:"type"`
}

type Pokemon struct {
	Name       string      `json:"name"`
	Experience int64       `json:"base_experience"`
	Height     int         `json:"height"`
	Weight     int         `json:"weight"`
	Stats      []Statistic `json:"stats"`
	Types      []Type      `json:"types"`
}

type PokemonAPI struct {
	name string
}

func NewPokemonAPI() *PokemonAPI {
	return &PokemonAPI{}
}

func (p *PokemonAPI) Path() string {
	return PokeAPIBaseURL + "/pokemon/" + p.name
}

func (p *PokemonAPI) Update(name string) {
	p.name = name
}

func (p *PokemonAPI) Get(url string) (any, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected 200 OK got %s", res.Status)
	}

	var pokemon Pokemon
	if err := json.NewDecoder(res.Body).Decode(&pokemon); err != nil {
		return nil, fmt.Errorf("decoding location area: %w", err)
	}

	return pokemon, nil
}

func (p *PokemonAPI) Print(pokemon Pokemon) {
	pokemonDetails := fmt.Sprintf(`Name: %s
Height: %d
Weight: %d
Stats:
`, pokemon.Name, pokemon.Height, pokemon.Weight)

	for _, stat := range pokemon.Stats {
		pokemonDetails += fmt.Sprintf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	pokemonDetails += "Types:\n"
	for _, t := range pokemon.Types {
		pokemonDetails += fmt.Sprintf("  - %s\n", t.Type.Name)
	}
	fmt.Print(pokemonDetails)
}
