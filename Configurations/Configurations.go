package configurations

import (
	"embed"

	"gopkg.in/yaml.v3"
)

//go:embed config.yml
var configData embed.FS

var config conf

type conf struct {
	Steam_api_key string `yaml:"steam_api_key"`
	Teams_folder  string `yaml:"teams_folder"`
	Match_folder  string `yaml:"match_folder"`
}

func load() {
	data, _ := configData.ReadFile("config.yml")
	yaml.Unmarshal(data, &config)
}

func GetConf() conf {
	if config.Steam_api_key == "" {
		load()
	}
	if config.Steam_api_key == "" {
		panic("No steam api key found")
	}
	return config
}
