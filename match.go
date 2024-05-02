package main

import (
	"flag"
	"fmt"
	"slices"

	o "glw.dk/csmatcher/Objects"
	utils "glw.dk/csmatcher/Utils"
)

var (
	execute = flag.String("execute", "", "What to execute, create match, or team")
)

func main() {
	flag.Parse()
	fmt.Println("CS:GO Match creator")
	if *execute == "" {
		fmt.Println("No command given")
		return
	}

	fmt.Println("Executing: ", *execute)
	switch *execute {
	case "match":
		o.CreateMatch()

	case "team":
		fmt.Println("Creating team:")
		var teamName string
		players := make([]string, 0)
		for {
			teamName = utils.InputPrompt("Team name: ")
			if ok := o.TeamExists(teamName); ok {
				fmt.Println("Team already exists")
				continue
			} else {
				break
			}
		}
		for {
			playerUrl := utils.InputPrompt("Player steam url/id (empty exit): ")
			if playerUrl == "" {
				break
			}
			if slices.Contains(players, playerUrl) {
				fmt.Println("Player already added")
				continue
			}
			players = append(players, playerUrl)
		}

		fmt.Printf("Team name: %s\n", teamName)
		fmt.Printf("Players: %v\n", players)

		o.CreateTeam(teamName, players)

	default:
		fmt.Println("Unknown command")
	}
}
