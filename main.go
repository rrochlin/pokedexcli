package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	// need to add in the help command during runtime to avoid dependancy loop.
	conf := config{}
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
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

func commandExit(conf *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf *config) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n")
	for _, cmd := range getCommands() {
		fmt.Println(fmt.Sprintf("%v: %v", cmd.name, cmd.description))
	}
	return nil
}

type config struct {
	Next     *string
	Previous *string
}

type LocationArea struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func commandMap(conf *config) error {
	var url string
	if conf.Next == nil {
		url = "https://pokeapi.co/api/v2/location-area/"
	} else {
		url = *conf.Next
	}

	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("could not get location data: %v", err)
	}

	defer res.Body.Close()
	if res.StatusCode > 299 {
		return fmt.Errorf("Server did not respond with OK, StatusCode:%v", res.StatusCode)
	}

	var locationSearch LocationArea
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Could not read response with io: %v", err)
	}

	if err = json.Unmarshal(body, &locationSearch); err != nil {
		return fmt.Errorf("Unmarshalling response failed: %v", err)
	}

	conf.Next = locationSearch.Next
	conf.Previous = locationSearch.Previous
	for _, loc := range locationSearch.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandMapB(conf *config) error {
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
