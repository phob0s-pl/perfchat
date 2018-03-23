package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/phob0s-pl/perfchat/chat"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	WithPrefix  bool
}

type API struct {
	Engine *chat.Chat
}

func NewAPI() *API {
	return &API{
		Engine: chat.NewChat(),
	}
}

// getRequestingUser returns requesting user
func (a *API) getRequestingUser(r *http.Request) (*chat.User, error) {
	id, token, _ := r.BasicAuth()
	user, err := a.Engine.GetUserByAuth(id, token)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUsers returns all users in chat
func (a *API) GetUsers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if _, ok := a.isUser(w, r); !ok {
		return
	}

	var users []User
	engineUsers := a.Engine.ListUsers()
	for _, engineUser := range engineUsers {
		users = append(users, User{Name: engineUser.Name, Role: engineUser.Role})
	}

	payload, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(payload)
}

// AddUser adds user to chat
// Only admin can add user
func (a *API) AddUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		user = &User{}
	)

	if _, ok := a.isUser(w, r); !ok {
		return
	}

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(content, user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := a.Engine.AddUser(
		&chat.User{
			AuthID: user.AuthID,
			Name:   user.Name,
			Role:   user.Role,
			Token:  user.Token,
		}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// Ping respongs wint 200 to ping message
func (a *API) Ping(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
}

// GetRooms returns list of rooms in chat
func (a *API) GetRooms(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if _, ok := a.isUser(w, r); !ok {
		return
	}

	var rooms []Room
	engineRooms := a.Engine.ListRooms()
	for _, engineRoom := range engineRooms {
		var userlist []string
		for _, userInRoom := range engineRoom.Users {
			userlist = append(userlist, userInRoom.Name)
		}
		rooms = append(rooms, Room{
			Name:    engineRoom.Name,
			Creator: engineRoom.Creator,
			Users:   userlist})
	}

	payload, err := json.Marshal(rooms)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(payload)
}

// CreateRoom creates new room and sets user as owner
func (a *API) CreateRoom(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	user, ok := a.isUser(w, r)
	if !ok {
		return
	}

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	room := &Room{}
	if err := json.Unmarshal(content, room); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := a.Engine.AddRoom(&chat.Room{
		Creator: user.Name,
		Name:    room.Name,
		Users:   []*chat.User{user},
	}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}

// RoomDelete deletes room from server
// note: user must be owner of a room
func (a *API) RoomDelete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if _, ok := a.isUser(w, r); !ok {
		return
	}
}

// RoomJoin joins user to room
func (a *API) RoomJoin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	user, ok := a.isUser(w, r)
	if !ok {
		return
	}

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	room := &Room{}
	if err := json.Unmarshal(content, room); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := a.Engine.JoinRoom(user.Name, room.Name); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// RoomExit exits user to room
// note: user can't be owner of the room
func (a *API) RoomExit(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	user, ok := a.isUser(w, r)
	if !ok {
		return
	}

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	room := &Room{}
	if err := json.Unmarshal(content, room); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := a.Engine.ExitRoom(user.Name, room.Name); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// Checks if request was done by user
// if not sets StatusUnauthorized on response and returns false
func (a *API) isUser(w http.ResponseWriter, r *http.Request) (*chat.User, bool) {
	user, err := a.getRequestingUser(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return nil, false
	}
	return user, true
}
