package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
)

// Config for this project
type Config struct {
	DB             string
	TelegramChatID string
	TelegramToken  string
}

// Default config
func (c *Config) Default() error {
	if c.DB == "" {
		return errors.New("please config postgres db (config.json)")
	}

	return nil
}

func (c *Config) Save() error {
	log.Println("CONFIG: save config.json")
	b, _ := json.MarshalIndent(c, "", "  ")
	return ioutil.WriteFile("config.json", b, 0644)
}

// Parse config from json byte
func Parse(b []byte, path string) (*Config, error) {
	c := &Config{}
	if err := json.Unmarshal(b, c); err != nil {
		return nil, err
	}
	if err := c.Default(); err != nil {
		return nil, err
	}

	return c, nil
}

// Read the config from `path`
func Read(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(b, path)
}

// Create the sample config file
func Create() error {
	log.Println("CONFIG: create config.json")
	b, _ := json.MarshalIndent(Config{}, "", "  ")
	return ioutil.WriteFile("config.json", b, 0644)
}
