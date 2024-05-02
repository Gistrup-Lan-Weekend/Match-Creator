package objects

import (
	"encoding/json"
	"fmt"
	"os"

	c "glw.dk/csmatcher/Configurations"
	u "glw.dk/csmatcher/Utils"
)

type Match struct {
	MatchId         string              `json:"matchId"`
	Team1           Team                `json:"team1"`
	Team2           Team                `json:"team2"`
	Num_maps        int                 `json:"num_maps"`
	Maplist         []string            `json:"maplist"`
	Players_pr_team int                 `json:"players_pr_team"`
	Cvars           []map[string]string `json:"cvars"`
}

func CreateMatch() (m Match) {
	os.MkdirAll(c.GetConf().Match_folder, os.ModePerm)
	m = Match{}
	// Select team
	teams := ReadTeams()
	firstTeamName := u.InputPrompt("Select 1'st team name: ")
	secondTeamName := u.InputPrompt("Select 2'nd team name: ")
	for {

		for {
			if firstTeamName == secondTeamName {
				fmt.Println("Teams must be different")
				secondTeamName = u.InputPrompt("Select 2'nd team name: ")
			} else {
				break
			}
		}

		for _, team := range teams {
			if team.TeamName == firstTeamName {
				m.Team1 = team
			}
			if team.TeamName == secondTeamName {
				m.Team2 = team
			}
		}
		if m.Team1.TeamName != "" && m.Team2.TeamName != "" {
			break
		} else {
			if m.Team1.TeamName == "" {
				firstTeamName = u.InputPrompt("Select 1'st team name: ")
			} else {
				secondTeamName = u.InputPrompt("Select 2'nd team name: ")
			}
		}
	}
	for {
		m.Num_maps = u.InputPromptInt("Number of maps: ")
		if m.Num_maps < 1 {
			fmt.Println("Number of maps must be greater than 0")
		} else if m.Num_maps > 5 {
			fmt.Println("Number of maps must be less than 6")
		} else {
			break
		}
	}

	for {
		m.Players_pr_team = u.InputPromptInt("Number of players pr team: ")
		if m.Players_pr_team < 1 {
			fmt.Println("Number of players pr team must be greater than 0")
		} else if m.Players_pr_team > 5 {
			fmt.Println("Number of players pr team must be less than 6")
		} else {
			break
		}
	}

	fmt.Printf("Maps:")
	for k, v := range Maps {
		fmt.Printf("%s: %s\n", k, v)
	}
	for {
		sel := u.InputPrompt("Select map: ")
		if selMap, ok := Maps[sel]; ok {
			m.Maplist = append(m.Maplist, selMap)
		} else {
			fmt.Printf("Invalid map")
		}
		if len(m.Maplist) == m.Num_maps {
			break
		}
	}
	file, err := os.Create(c.GetConf().Match_folder + firstTeamName + "_vs_" + secondTeamName + ".json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// Marshal the data
	data, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	// Write the data
	file.Write(data)
	return
}

func (m Match) GetMatchAsJson() ([]byte, error) {
	return json.Marshal(m)
}
