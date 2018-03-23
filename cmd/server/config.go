package main

import (
	"github.com/BurntSushi/toml"
)

// Config represents server configuration
type Config struct {
	// Address to listen to
	Address string

	// Debug if set to true will enable debug logs
	Debug bool

	// AuthID is administrator AuthID
	AuthID string

	//AdminPass is administrator Token
	Token string
}

// ReadConfig reads config from file
func ReadConfig(path string) (*Config, error) {
	cfg := &Config{}
	_, err := toml.DecodeFile(path, cfg)
	return cfg, err
}
