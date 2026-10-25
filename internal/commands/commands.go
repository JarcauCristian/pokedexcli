// Package commands provides the implementation for all the commands that can be called by the CLI.
package commands

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"

	"github.com/JarcauCristian/pokedexcli/internal/api"
)

type CliCommand struct {
	Name        string
	Description string
	Callback    func(config *Config, args ...string) error
}

func CommandRegistry() map[string]CliCommand {
	commands := map[string]CliCommand{
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			Callback:    exitCommand,
		},
		"help": {
			Name:        "help",
			Description: "Displays a help message",
			Callback:    helpCommand,
		},
		"map": {
			Name:        "map",
			Description: "Returns the next 20 Pokemon world areas until all are printed",
			Callback:    mapCommand,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Same as map, but returns the previous 20 entries instead of the next",
			Callback:    mapbCommand,
		},
		"explore": {
			Name:        "explore",
			Description: "Explores an area in more details",
			Callback:    exploreCommand,
		},
		"catch": {
			Name:        "catch",
			Description: "Try and catch a pokemon specified by name",
			Callback:    catchCommand,
		},
		"pokedex": {
			Name:        "pokedex",
			Description: "Returns the list of pokemons that were caught by the user",
			Callback:    pokedexCommand,
		},
		"inspect": {
			Name:        "inspect",
			Description: "Inspects a specific pokemon specified by name",
			Callback:    inspectCommand,
		},
	}

	return commands
}

func exitCommand(config *Config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func helpCommand(config *Config, args ...string) error {
	helpText := `Welcome to the Pokedex!
Usage:

`
	for _, command := range config.Commands {
		helpText += fmt.Sprintf("%s: %s\n", command.Name, command.Description)
	}
	fmt.Println(helpText)
	return nil
}

func mapCommand(config *Config, args ...string) error {
	var response api.LocationAreaResponse
	var err error

	areasAPI, ok := config.Endpoints["areas"].(*api.LocationAreasAPI)
	if !ok {
		return fmt.Errorf("failed to convert the areas endpoint")
	}

	url := areasAPI.Path()
	if areasAPI.NextURL != nil {
		url = *areasAPI.NextURL
	}

	if value, ok := config.Cache.Get(url); ok {
		err = json.Unmarshal(value, &response)
	} else {
		res, err := areasAPI.Get(url)
		if err != nil {
			return err
		}
		response, ok = res.(api.LocationAreaResponse)
		if !ok {
			return fmt.Errorf("failed to parse the response")
		}
	}

	areasAPI.Update(response)
	for _, result := range response.Results {
		fmt.Println(result.Name)
	}

	value, err := json.Marshal(response)
	if err != nil {
		return err
	}
	config.Cache.Add(url, value)

	return nil
}

func mapbCommand(config *Config, args ...string) error {
	areasAPI, ok := config.Endpoints["areas"].(*api.LocationAreasAPI)
	if !ok {
		return fmt.Errorf("failed to convert the areas endpoint")
	}

	url := areasAPI.Path()
	if areasAPI.PrevURL != nil {
		url = *areasAPI.PrevURL
	}

	var (
		response api.LocationAreaResponse
		err      error
	)
	if value, ok := config.Cache.Get(url); ok {
		err = json.Unmarshal(value, &response)
		if err != nil {
			return err
		}
	} else {
		res, err := areasAPI.Get(url)
		if err != nil {
			return err
		}

		response, ok = res.(api.LocationAreaResponse)
		if !ok {
			return fmt.Errorf("failed to parse the response")
		}
	}

	areasAPI.Update(response)
	for _, result := range response.Results {
		fmt.Println(result.Name)
	}

	value, err := json.Marshal(response)
	if err != nil {
		return err
	}
	config.Cache.Add(url, value)

	return nil
}

func exploreCommand(config *Config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("this function accepts only a single name argument, received %d", len(args))
	}

	areaAPI, ok := config.Endpoints["area"].(*api.LocationAreaAPI)
	if !ok {
		return fmt.Errorf("failed to convert the areas endpoint")
	}
	areaAPI.Update(args[0])

	url := areaAPI.Path()

	fmt.Println(
		fmt.Sprintf("Exploring %s...", args[0]),
	)

	var (
		response []string
		err      error
	)
	if value, ok := config.Cache.Get(url); ok {
		err = json.Unmarshal(value, response)
		if err != nil {
			return err
		}
	} else {
		res, err := areaAPI.Get(url)
		if err != nil {
			return err
		}

		response, ok = res.([]string)
		if !ok {
			return fmt.Errorf("failed to parse the response")
		}
	}

	fmt.Println("Found Pokemon:")
	for _, name := range response {
		fmt.Printf("- %s\n", name)
	}
	return nil
}

func catchCommand(config *Config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("this function accepts only a single name argument, received %d", len(args))
	}
	name := args[0]

	pokemonAPI, ok := config.Endpoints["pokemon"].(*api.PokemonAPI)
	if !ok {
		return fmt.Errorf("failed to convert the areas endpoint")
	}
	pokemonAPI.Update(name)

	var (
		pokemon api.Pokemon
		err     error
	)
	url := pokemonAPI.Path()
	if value, ok := config.Cache.Get(url); ok {
		err = json.Unmarshal(value, &pokemon)
		if err != nil {
			return err
		}
	} else {
		res, err := pokemonAPI.Get(url)
		if err != nil {
			return err
		}

		pokemon, ok = res.(api.Pokemon)
		if !ok {
			return fmt.Errorf("failed to parse the response")
		}
	}

	value, err := json.Marshal(pokemon)
	if err != nil {
		return err
	}
	config.Cache.Add(url, value)

	fmt.Printf("Throwing a Pokeball at %s...\n", name)

	r := rand.New(rand.NewSource(0))
	chance := r.Intn(int(pokemon.Experience))

	if chance <= int(pokemon.Experience)/2 {
		fmt.Printf("%s was caught!\n", name)
		config.Pokedex.Add(pokemon)
	} else {
		fmt.Printf("%s escaped!\n", name)
	}

	return nil
}

func pokedexCommand(config *Config, args ...string) error {
	config.Pokedex.List()
	return nil
}

func inspectCommand(config *Config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("this function accepts only a single name argument, received %d", len(args))
	}
	name := args[0]

	pokemonAPI, ok := config.Endpoints["pokemon"].(*api.PokemonAPI)
	if !ok {
		return fmt.Errorf("failed to convert the areas endpoint")
	}
	pokemonAPI.Update(name)

	var (
		pokemon api.Pokemon
		err     error
	)
	url := pokemonAPI.Path()
	if value, ok := config.Cache.Get(url); ok {
		err = json.Unmarshal(value, &pokemon)
		if err != nil {
			return err
		}
	} else {
		res, err := pokemonAPI.Get(url)
		if err != nil {
			return err
		}

		pokemon, ok = res.(api.Pokemon)
		if !ok {
			return fmt.Errorf("failed to parse the response")
		}
	}

	value, err := json.Marshal(pokemon)
	if err != nil {
		return err
	}
	config.Cache.Add(url, value)

	pokemonAPI.Print(pokemon)

	return nil
}
