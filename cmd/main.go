package main

import (
	"time"

	"github.com/JarcauCristian/pokedexcli/internal/api"
	"github.com/JarcauCristian/pokedexcli/internal/commands"
	"github.com/JarcauCristian/pokedexcli/internal/pokecache"
	"github.com/JarcauCristian/pokedexcli/internal/pokedex"
	"github.com/JarcauCristian/pokedexcli/internal/repl"
)

func main() {
	c := commands.CommandRegistry()
	pokedex := pokedex.NewPokedex()
	cache := pokecache.NewCache(5 * time.Second)
	defer cache.Stop()

	config := commands.Config{
		Commands:  c,
		Endpoints: api.EndpointsRegistry(),
		Cache:     cache,
		Pokedex:   pokedex,
	}
	repl.REPL(&config)
}
