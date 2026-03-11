package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

const DefaultEndpoint = "https://homebox.mizobuchi.dev"

type Config struct {
	Endpoint string `json:"endpoint,omitempty"`
	Token    string `json:"token,omitempty"`
}

func configPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "hb", "config.json")
}

func Load() Config {
	data, err := os.ReadFile(configPath())
	if err != nil {
		return Config{}
	}
	var c Config
	_ = json.Unmarshal(data, &c)
	return c
}

func Save(c Config) error {
	p := configPath()
	if err := os.MkdirAll(filepath.Dir(p), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0600)
}

func GetEndpoint() string {
	if e := os.Getenv("HB_ENDPOINT"); e != "" {
		return strings.TrimRight(e, "/")
	}
	c := Load()
	if c.Endpoint != "" {
		return strings.TrimRight(c.Endpoint, "/")
	}
	return DefaultEndpoint
}

func GetToken() string {
	if t := os.Getenv("HB_TOKEN"); t != "" {
		return t
	}
	return Load().Token
}

func SetToken(token string) error {
	c := Load()
	c.Token = token
	return Save(c)
}

func ClearToken() error {
	c := Load()
	c.Token = ""
	return Save(c)
}

func SetEndpoint(endpoint string) error {
	c := Load()
	c.Endpoint = strings.TrimRight(endpoint, "/")
	return Save(c)
}
