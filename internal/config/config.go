package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"strings"

	"gps/internal/gpt"
)

const DefaultPath string = "~/.gps/config.json"

type Config struct {
	Gpt gpt.GptConfig `json:"gpt"`
}

func LoadDefaultConfig() (*Config, error) {
	return LoadConfig(DefaultPath)
}

// LoadConfig reads a YAML configuration file and returns the parsed configuration.
func LoadConfig(filename string) (*Config, error) {
	if filename == "" {
		filename = DefaultPath
	}

	if strings.HasPrefix(filename, "~") {
		usr, err := user.Current()
		if err != nil {
			return nil, err
		}

		homeDir := usr.HomeDir
		filename = strings.Replace(filename, "~", homeDir, 1)
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, err
	}

	// read the file contents
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// parse the YAML contents
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
