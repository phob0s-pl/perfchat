package main

import (
	"github.com/BurntSushi/toml"
)

// Config represents client configuration
type Config struct {
	// Address to listen to
	Address string

	// Debug if set to true will enable debug logs
	Debug bool

	// AuthID is administrator AuthID
	AuthID string

	//AdminPass is administrator Token
	Token string

	//Workers is count of workers simulating real users
	Workers uint

	// RoomOp is time in ms how often room operation is dome
	RoomOp uint

	// MessageChance is time in ms is chance to join room with random user and send message
	MessageToUserChance uint
}

// ReadConfig reads config from file
func ReadConfig(path string) (*Config, error) {
	cfg := &Config{}
	_, err := toml.DecodeFile(path, cfg)
	return cfg, err
}
