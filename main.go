package main

import (
	"bufio"
	"fmt"
	"github.com/rrochlin/pokedexcli/internals"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	conf := internals.Config{
		Cache:   internals.NewCache(5 * time.Second),
		Pokedex: make(map[string]internals.Pokemon),
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		command := cleanInput(scanner.Text())
		if len(command) < 1 {
			fmt.Println("Invalid Command")
			continue
		}
		if cmd, ok := getCommands()[command[0]]; ok {
			if err := cmd.callback(&conf, command[1:]...); err != nil {
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
		"catch": {
			name:        "catch",
			description: "Attempt to catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a pokemon from the pokedex",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List out all captured pokemon",
			callback:    commandPokedex,
		},
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func(*internals.Config, ...string) error
}

func commandExit(conf *internals.Config, _ ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf *internals.Config, _ ...string) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, cmd := range getCommands() {
		fmt.Println(fmt.Sprintf("%v: %v", cmd.name, cmd.description))
	}
	return nil
}

func commandMap(conf *internals.Config, _ ...string) error {
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

func commandMapB(conf *internals.Config, _ ...string) error {
	if conf.Previous == nil {
		return fmt.Errorf("No previous locations to show")
	}
	t := conf.Next
	conf.Next = conf.Previous
	conf.Previous = t
	if err := commandMap(conf, ""); err != nil {
		return err
	}
	return nil
}

func commandExplore(conf *internals.Config, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("Not enough args supplied for explore")
	}
	area := args[0]
	better_name := strings.ReplaceAll(area, "-", " ")
	fmt.Printf("Exploring %v...\n", better_name)
	pokeArea, err := internals.GetPokeArea(conf, area)
	if err != nil {
		return err
	}
	fmt.Println("Found Pokemon:")
	for _, encounter := range pokeArea.PokemonEncounters {
		fmt.Printf(" - %v\n", encounter.Pokemon.Name)
	}
	return nil
}

func commandCatch(conf *internals.Config, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("Not enough args supplied for catch")
	}
	pokemonSearch := args[0]
	pokemon, err := internals.GetPokemon(conf, pokemonSearch)
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %v...\n", pokemon.Name)
	thresh := 0.25 + float64(3.0*pokemon.BaseExperience)/float64(2000.0+5.0*pokemon.BaseExperience)
	fmt.Println(thresh)
	if rand.Float64() > thresh {
		conf.Pokedex[pokemon.Name] = pokemon
		fmt.Printf("%v was caught!\n", pokemon.Name)
	} else {
		fmt.Printf("%v escaped!\n", pokemon.Name)
	}
	return nil
}

func commandInspect(conf *internals.Config, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("Not enough args supplied for inspect")
	}
	if mon, ok := conf.Pokedex[args[0]]; !ok {
		fmt.Println("Pokemon has not been caught")
	} else {
		fmt.Printf("Name: %v\n", mon.Name)
		fmt.Printf("Height: %v\n", mon.Height)
		fmt.Printf("Weight: %v\n", mon.Weight)
		fmt.Printf("Stats: \n")
		for _, stat := range mon.Stats {
			fmt.Printf(" -%v: %v\n", stat.Stat.Name, stat.BaseStat)
		}
		for _, tp := range mon.Types {
			fmt.Printf(" - %v\n", tp.Type.Name)
		}
	}
	return nil
}

func commandPokedex(conf *internals.Config, _ ...string) error {
	fmt.Println("Your Pokedex:")
	for name := range conf.Pokedex {
		fmt.Printf(" - %v\n", name)
	}
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
