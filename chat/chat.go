package chat

import (
	"errors"
	"sync"
)

var (
	// ErrNotFound is returned when resource is not found
	ErrNotFound = errors.New("not found")
	// ErrExists is returned when resource already exists
	ErrExists = errors.New("already exists")
	// ErrNoResources is returned when limits are reached
	ErrNoResources = errors.New("out of resources")
	// ErrMissingArg is returned when some data is missing
	ErrMissingArg = errors.New("missing argument")
	// ErrNotPermit is returned when operation is not permited
	ErrNotPermit = errors.New("not permited")
)

const (
	defaultUsersLimit = 1000
	defaultRoomsLimit = 1000
	mainRoom          = "main"
)

// Chat is main chat engine
type Chat struct {
	roomLimit  uint
	usersLimit uint
	sync.Mutex
	users map[string]*User
	rooms map[string]*Room
}

// SetRoomLimit sets room limit to non default value
func (c *Chat) SetRoomLimit(limit uint) {
	c.roomLimit = limit
}

// SetUsersLimit sets users limit to non default value
func (c *Chat) SetUsersLimit(limit uint) {
	c.usersLimit = limit
}

// NewChat returns new chat
func NewChat() *Chat {
	c := &Chat{
		users:      make(map[string]*User),
		rooms:      make(map[string]*Room),
		roomLimit:  defaultRoomsLimit,
		usersLimit: defaultUsersLimit,
	}
	c.rooms[mainRoom] = &Room{Name: mainRoom}
	return c
}

// AddUser adds user to chat
func (c *Chat) AddUser(user *User) error {
	c.Lock()
	defer c.Unlock()

	if uint(c.UsersCount()) >= c.usersLimit {
		return ErrNoResources
	}

	if _, ok := c.users[user.Name]; ok {
		return ErrExists
	}

	c.users[user.Name] = user
	return nil
}

// AddRoom adds user to chat
func (c *Chat) AddRoom(room *Room) error {
	c.Lock()
	defer c.Unlock()

	if len(room.Users) != 1 || room.Name == "" || room.Creator == "" {
		return ErrMissingArg
	}

	if _, ok := c.rooms[room.Name]; ok {
		return ErrExists
	}

	c.rooms[room.Name] = room
	return nil
}

// JoinRoom joins user with name to room
func (c *Chat) JoinRoom(username, roomname string) error {
	c.Lock()
	defer c.Unlock()

	user, err := c.GetUserByName(username)
	if err != nil {
		return err
	}

	room, err := c.GetRoomByName(roomname)
	if err != nil {
		return err
	}

	return room.Join(user)
}

// ExitRoom removes user with name from room
func (c *Chat) ExitRoom(username, roomname string) error {
	c.Lock()
	defer c.Unlock()

	user, err := c.GetUserByName(username)
	if err != nil {
		return err
	}

	room, err := c.GetRoomByName(roomname)
	if err != nil {
		return err
	}

	return room.Exit(user)
}

// RoomExists checks if room exists
func (c *Chat) RoomExists(roomname string) bool {
	c.Lock()
	defer c.Unlock()
	_, ok := c.rooms[roomname]
	return ok
}

// UserExists checks if user with name exists
func (c *Chat) UserExists(username string) bool {
	c.Lock()
	defer c.Unlock()
	_, ok := c.users[username]
	return ok
}

// GetUserByName returns user by his name
func (c *Chat) GetUserByName(name string) (*User, error) {
	c.Lock()
	defer c.Unlock()
	user, ok := c.users[name]
	if !ok {
		return nil, ErrNotFound
	}
	return user, nil
}

// GetRoomByName returns room by its name
func (c *Chat) GetRoomByName(name string) (*Room, error) {
	c.Lock()
	defer c.Unlock()
	room, ok := c.rooms[name]
	if !ok {
		return nil, ErrNotFound
	}
	return room, nil
}

// ListUsers returns all users
func (c *Chat) ListUsers() (u []*User) {
	c.Lock()
	defer c.Unlock()
	for _, user := range c.users {
		u = append(u, user)
	}
	return u
}

// ListRooms returns all rooms
func (c *Chat) ListRooms() (r []*Room) {
	c.Lock()
	defer c.Unlock()
	for _, room := range c.rooms {
		r = append(r, room)
	}
	return r
}

// GetUserByAuth returns user by auth parameters
func (c *Chat) GetUserByAuth(id, token string) (*User, error) {
	for _, user := range c.users {
		if user.AuthID == id && user.Token == token {
			return user, nil
		}
	}
	return nil, ErrNotFound
}

// UsersCount returns number of users in chat
func (c *Chat) UsersCount() int {
	return len(c.users)
}

// RoomsCount returns number of rooms in chat
func (c *Chat) RoomsCount() int {
	return len(c.rooms)
}

// DeleteRoom deletes room from chat if exists
// and username is room creator
func (c *Chat) DeleteRoom(username, roomname string) error {
	c.Lock()
	defer c.Unlock()

	room, ok := c.rooms[roomname]
	if !ok {
		return ErrNotFound
	}

	if room.Creator != username {
		return ErrNotPermit
	}

	delete(c.rooms, roomname)
	return nil
}
