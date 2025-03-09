package main

import (
	"bufio"
	"fmt"
	"github.com/rrochlin/pokedexcli/internals"
	"os"
	"strings"
	"time"
)

func main() {
	conf := internals.Config{Cache: internals.NewCache(5 * time.Second)}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		command := cleanInput(scanner.Text())[0]
		if cmd, ok := getCommands()[command]; ok {
			if err := cmd.callback(&conf); err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}

	}
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Diplays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Display Next 20 Locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display Previous 20 Locations",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "Check for Pokemon located in an area",
			callback:    commandExplore,
		},
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func(*internals.Config) error
}

func commandExit(conf *internals.Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf *internals.Config) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, cmd := range getCommands() {
		fmt.Println(fmt.Sprintf("%v: %v", cmd.name, cmd.description))
	}
	return nil
}

func commandMap(conf *internals.Config) error {
	locationSearch, err := internals.GetPokeLocations(conf)
	if err != nil {
		return err
	}

	conf.Next = locationSearch.Next
	conf.Previous = locationSearch.Previous
	for _, loc := range locationSearch.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandMapB(conf *internals.Config) error {
	if conf.Previous == nil {
		return fmt.Errorf("No previous locations to show")
	}
	t := conf.Next
	conf.Next = conf.Previous
	conf.Previous = t
	if err := commandMap(conf); err != nil {
		return err
	}
	return nil
}

func commandExplore(conf *internals.Config) error {
	return nil
}

func cleanInput(text string) []string {
	processed := strings.ToLower(text)
	substrings := strings.Split(processed, " ")
	res := make([]string, 0)
	for _, str := range substrings {
		if str == "" {
			continue
		}
		res = append(res, str)
	}
	return res
}
