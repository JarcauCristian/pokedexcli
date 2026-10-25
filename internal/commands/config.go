package commands

import (
	"github.com/JarcauCristian/pokedexcli/internal/api"
	"github.com/JarcauCristian/pokedexcli/internal/pokecache"
	"github.com/JarcauCristian/pokedexcli/internal/pokedex"
)

type Config struct {
	Commands  map[string]CliCommand
	Endpoints map[string]api.Endpoint
	Cache     *pokecache.Cache
	Pokedex   *pokedex.Pokedex
}
