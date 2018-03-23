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

	if c.UserExists(user.Name) {
		return ErrExists
	}

	c.users[user.Name] = user
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

	return room.join(user)
}

// RoomExists checks if room exists
func (c *Chat) RoomExists(roomname string) bool {
	_, ok := c.rooms[roomname]
	return ok
}

// UserExists checks if user with name exists
func (c *Chat) UserExists(username string) bool {
	_, ok := c.users[username]
	return ok
}

// GetUserByName returns user by his name
func (c *Chat) GetUserByName(name string) (*User, error) {
	user, ok := c.users[name]
	if !ok {
		return nil, ErrNotFound
	}
	return user, nil
}

// GetRoomByName returns room by its name
func (c *Chat) GetRoomByName(name string) (*Room, error) {
	room, ok := c.rooms[name]
	if !ok {
		return nil, ErrNotFound
	}
	return room, nil
}

// ListUsers returns all users
func (c *Chat) ListUsers() (u []*User) {
	for _, user := range c.users {
		u = append(u, user)
	}
	return u
}

// ListRooms returns all rooms
func (c *Chat) ListRooms() (r []*Room) {
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
