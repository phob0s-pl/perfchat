package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/websocket"
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
	Engine           *chat.Chat
	upgrader         *websocket.Upgrader
	websocketClients map[string]*websocketClient
	message          chan *Message
	msgBuffer        map[string]chan *Message
}

func NewAPI() *API {
	return &API{
		Engine: chat.NewChat(),
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024 * 1024,
			WriteBufferSize: 1024 * 1024,
		},
		websocketClients: make(map[string]*websocketClient),
		message:          make(chan *Message),
		msgBuffer:        make(map[string]chan *Message),
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
	a.msgBuffer[user.Name] = make(chan *Message, 256)
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

// ReceiveMessage receives message from client
func (a *API) ReceiveMessage(w http.ResponseWriter, r *http.Request) {
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
	msg := &Message{}
	if err := json.Unmarshal(content, msg); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	msg.User = user.Name

	room, err := a.Engine.GetRoomByName(msg.Room)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, roomUser := range room.Users {
		c, ok := a.msgBuffer[roomUser.Name]
		if !ok {
			continue
		}
		if len(c) < cap(c) {
			c <- msg
		}
	}
}

// SendMessage sends messages to client
func (a *API) SendMessage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	user, ok := a.isUser(w, r)
	if !ok {
		return
	}

	channel, ok := a.msgBuffer[user.Name]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	msgs := unchannelMsg(channel)

	payload, err := json.Marshal(msgs)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(payload)
}

func unchannelMsg(channel chan *Message) (msgs []Message) {
	for {
		select {
		case msg := <-channel:
			msgs = append(msgs, *msg)
		default:
			return
		}
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

// Websocket handles messages from clients
func (a *API) Websocket(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	user, ok := a.isUser(w, r)
	if !ok {
		return
	}
	conn, err := a.upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	client := &websocketClient{
		conn: conn,
		msg:  make(chan *Message, 256),
		name: user.Name,
	}

	a.websocketClients[user.Name] = client
	a.writeClientMessage(client)
	a.readClientMessage(client)
}
