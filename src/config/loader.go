package config

import (
	"bytes"
	"errors"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func LoadFromData(data string) error {
	decoded, err := Decode(data)
	if err != nil {
		return err
	}
	
	viper.SetConfigType("ini")
	return viper.ReadConfig(bytes.NewReader(decoded))
}

func LoadFromEnv() (bool, error) {
	data := os.Getenv("WIRETAP_CONFIG_DATA")
	if data == "" {
		return false, nil
	}
	
	err := LoadFromData(data)
	return true, err
}

func LoadFromFile(path string) (bool, error) {
	if !strings.HasSuffix(path, ".enc") {
		return false, nil
	}
	
	content, err := os.ReadFile(path)
	if err != nil {
		return true, err
	}
	
	data := strings.TrimSpace(string(content))
	err = LoadFromData(data)
	return true, err
}

func LoadFromURL(cfg FetchConfig) error {
	data, err := FetchRemote(cfg)
	if err != nil {
		return err
	}
	
	dataStr := strings.TrimSpace(string(data))
	return LoadFromData(dataStr)
}

func TryLoad(configData string, configFile string, configURL string, authHeader string, insecure bool) error {
	if configData != "" {
		return LoadFromData(configData)
	}
	
	if configURL != "" {
		return LoadFromURL(FetchConfig{
			URL:        configURL,
			AuthHeader: authHeader,
			Insecure:   insecure,
		})
	}
	
	if loaded, err := LoadFromEnv(); loaded {
		return err
	}
	
	if configFile != "" {
		if loaded, err := LoadFromFile(configFile); loaded {
			return err
		}
	}
	
	return errors.New("no config")
}
