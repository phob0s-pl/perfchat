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
)

// Chat is main chat engine
type Chat struct {
	sync.Mutex
	users map[string]*User
	rooms map[string]*Room
}

// NewChat returns new chat
func NewChat() *Chat {
	return &Chat{
		users: make(map[string]*User),
		rooms: make(map[string]*Room),
	}
}

// AddUser adds user to chat
func (c *Chat) AddUser(user *User) error {
	c.Lock()
	defer c.Unlock()

	if c.UserExists(user) {
		return ErrExists
	}

	c.users[user.Name] = user
	return nil
}

// UserExists checks if user exists
func (c *Chat) UserExists(user *User) bool {
	_, ok := c.users[user.Name]
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

// ListUsers returns all users
func (c *Chat) ListUsers() (u []*User) {
	for _, user := range c.users {
		u = append(u, user)
	}
	return u
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
