package objects

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"os"
	"strings"

	c "glw.dk/csmatcher/Configurations"
	"gopkg.in/yaml.v3"
)

type Team struct {
	TeamName string            `json:"name"`
	Players  map[string]string `json:"players"`
}

func ReadTeams() []Team {
	teams := make([]Team, 0)
	// Read all files in the folder
	files, err := os.ReadDir(c.GetConf().Teams_folder)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		teamName := strings.Split(file.Name(), ".")[0]
		// Read the file
		data, err := os.ReadFile(c.GetConf().Teams_folder + file.Name())
		if err != nil {
			panic(err)
		}
		// Unmarshal the data
		teamPlayers := make(map[string]string, 0)
		yaml.Unmarshal(data, &teamPlayers)
		// Create the team
		teams = append(teams, Team{TeamName: teamName, Players: teamPlayers})
	}
	return teams
}

func GetPlayersFromUrls(url string) (player map[string]string) {
	player = make(map[string]string)
	steamId := ""
	if strings.Contains(url, "profiles") {
		steamId = strings.Replace(url, "https://steamcommunity.com/profiles/", "", -1)
	} else {
		strippedUrl := strings.Replace(url, "https://steamcommunity.com/id/", "", -1)
		strippedUrl = strings.ReplaceAll(strippedUrl, "/", "")
		resolveSteamUrl := fmt.Sprintf("http://api.steampowered.com/ISteamUser/ResolveVanityURL/v0001/?key=%s&vanityurl=%s",
			c.GetConf().Steam_api_key, strippedUrl)
		resp, err := http.Get(resolveSteamUrl)
		if err != nil {
			fmt.Println("Error fetching steam data")
		} else {
			defer resp.Body.Close()
		}

		// Parse the steam data
		var steamData map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&steamData)
		steamId = steamData["response"].(map[string]interface{})["steamid"].(string)
	}
	// Get the player's nick
	getPlayerNick := fmt.Sprintf("http://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=%s&steamids=%s", c.GetConf().Steam_api_key, steamId)
	resp, err := http.Get(getPlayerNick)

	if err != nil {
		fmt.Println("Error fetching steam data")
		return nil
	}
	defer resp.Body.Close()

	// Parse the steam data
	var steamData map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&steamData)
	player[steamId] = steamData["response"].(map[string]interface{})["players"].([]interface{})[0].(map[string]interface{})["personaname"].(string)
	return
}

func CreateTeam(name string, playersUrlORID []string) {
	os.MkdirAll(c.GetConf().Teams_folder, os.ModePerm)

	players := make(map[string]string, 0)
	for _, urlOrID := range playersUrlORID {
		if strings.Contains(urlOrID, "steamcommunity.com") {
			maps.Copy(players, GetPlayersFromUrls(urlOrID))
		} else if strings.Contains(urlOrID, "765611") {
			maps.Copy(players, CreateTeamFromIds(urlOrID))
		} else {
			fmt.Println("Invalid url or id")
		}
	}
	file, err := os.Create(c.GetConf().Teams_folder + name + ".yml")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// Marshal the data
	data, err := yaml.Marshal(players)
	if err != nil {
		panic(err)
	}
	// Write the data
	file.Write(data)
}

func CreateTeamFromIds(id string) (player map[string]string) {
	player = make(map[string]string)
	// Get the player's nick
	getPlayerNick := fmt.Sprintf("http://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=%s&steamids=%s", c.GetConf().Steam_api_key, id)
	resp, err := http.Get(getPlayerNick)
	if err != nil {
		fmt.Println("Error fetching steam data")
		return nil
	}

	defer resp.Body.Close()

	var steamData map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&steamData)
	player[id] = steamData["response"].(map[string]interface{})["players"].([]interface{})[0].(map[string]interface{})["personaname"].(string)
	return
}

func TeamExists(teamName string) bool {
	files, err := os.ReadDir(c.GetConf().Teams_folder)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		name := strings.Split(file.Name(), ".")[0]
		if strings.ToLower(name) == strings.ToLower(teamName) {
			return true
		}
	}
	return false
}
