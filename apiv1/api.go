package api

import "fmt"

const (
	// Version is API version
	Version = "1"
	// RootPath ris root path to API calls
	RootPath = "perfchat"
)

const (
	// UsersCall is API call for:
	// - POST method adds user
	// - GET method lists all users
	UsersCall = "users"

	// PingCall for checking if API is online
	PingCall = "ping"
)

type User struct {
	Name   string `json:"name"`
	Role   string `json:"role"`
	AuthID string `json:"authid"`
	Token  string `json:"token"`
}

// GetPath returns path to API call
func GetPath(path string) string {
	return fmt.Sprintf("/%s/%s/%s/", RootPath, Version, path)
}
