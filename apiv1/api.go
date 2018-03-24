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

	// PingCall [GET] is for checking if API is online
	PingCall = "ping"

	// RoomsCall [GET] is for listing current rooms
	RoomsCall = "rooms"

	// RoomsCreateCall [POST] is for creating new room
	RoomsCreateCall = "rooms/create"

	// RoomsDeleteCall [POST] is for deleting  room
	RoomsDeleteCall = "rooms/delete"

	// RoomsJoinCall [POST] is for joining selected room
	RoomsJoinCall = "rooms/join"

	// RoomsExitCall [POST] is for exiting from room
	RoomsExitCall = "rooms/exit"

	// MessageCall [POST] posts single message [GET] retrieves all messages
	MessageCall = "message"

	// WsPath is websocket path
	WsPath = "ws"

	// StatsCall returns JSON statistics
	StatsCall = "stats"
)

// User is structure for manipulating user related calls
type User struct {
	Name   string `json:"name"`
	Role   string `json:"role"`
	AuthID string `json:"authid"`
	Token  string `json:"token"`
	// Rooms represents rooms user is currently joined in
	Rooms []string `json:"rooms"`
}

// Room is structure for manipulating room related calls
type Room struct {
	Name    string   `json:"name"`
	Creator string   `json:"creator"`
	Users   []string `json:"users"`
}

type Message struct {
	User    string `json:"user"`
	Room    string `json:"room"`
	Content string `json:"content"`
}

// GetPath returns path to API call
func GetPath(path string) string {
	return fmt.Sprintf("/%s/%s/%s/", RootPath, Version, path)
}
