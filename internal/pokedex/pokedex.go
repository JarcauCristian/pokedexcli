package pokedex

import (
	"fmt"
	"sync"

	"github.com/JarcauCristian/pokedexcli/internal/api"
)

type Pokedex struct {
	mu       sync.Mutex
	pokemons map[string]api.Pokemon
}

func NewPokedex() *Pokedex {
	return &Pokedex{
		pokemons: make(map[string]api.Pokemon),
	}
}

func (p *Pokedex) Add(pokemon api.Pokemon) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.pokemons[pokemon.Name] = pokemon
}

func (p *Pokedex) Get(name string) (api.Pokemon, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	pokemon, ok := p.pokemons[name]
	return pokemon, ok
}

func (p *Pokedex) List() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, pokemon := range p.pokemons {
		fmt.Printf(`- Name: %s
  Experience: %d
-------------
`, pokemon.Name, pokemon.Experience)
	}
}
