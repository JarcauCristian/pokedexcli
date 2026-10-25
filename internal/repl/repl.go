// Package repl provides the read-eval-print loop function for the CLI.
package repl

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/JarcauCristian/pokedexcli/internal/commands"
)

var invalidCharsRe = regexp.MustCompile(`[^a-zA-Z0-9_-]`)

func cleanInput(text string) []string {
	entries := strings.Split(
		strings.ToLower(
			strings.Trim(text, " "),
		),
		" ",
	)

	for i := range entries {
		entries[i] = invalidCharsRe.ReplaceAllString(entries[i], "")
	}
	return entries
}

func REPL(config *commands.Config) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\033[34mPokedex\033[0m > ")
		if !scanner.Scan() {
			continue
		}

		text := scanner.Text()
		input := cleanInput(text)
		if command, ok := config.Commands[input[0]]; ok {
			if err := command.Callback(config, input[1:]...); err != nil {
				errMsg := fmt.Sprintf("\033[31mError encountered while executing %s, err=%v\033[0m", command.Name, err)
				fmt.Println(errMsg)
			}
		} else {
			fmt.Println("\033[33mUnknown command\033[0m")
		}
	}
}
