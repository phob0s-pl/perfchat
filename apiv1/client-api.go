package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"io/ioutil"

	"github.com/phob0s-pl/perfchat/chat"
)

// Client is API client
type Client struct {
	user       *chat.User
	httpClient *http.Client
	serverAddr string
}

// NewClient returns API client
// user represents user performing operations
// addr is server address
func NewClient(user *chat.User, addr string) *Client {
	return &Client{
		httpClient: &http.Client{},
		user:       user,
		serverAddr: addr,
	}
}

func (c *Client) requestPath(apicall string) string {
	return fmt.Sprintf("http://%s%s", c.serverAddr, GetPath(apicall))
}

// AddUser adds new chat user
// Note: need to have admin priviliges
func (c *Client) AddUser(user *chat.User) error {
	payload, err := json.Marshal(&User{
		AuthID: user.AuthID,
		Name:   user.Name,
		Role:   user.Role,
		Token:  user.Token,
	})

	if err != nil {
		return fmt.Errorf("AddUser: %s", err)
	}

	url := c.requestPath(UsersCall) + user.Name
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("AddUser: %s", err)
	}
	request.SetBasicAuth(c.user.AuthID, c.user.Token)

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("AddUser: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("AddUser: %s", http.StatusText(resp.StatusCode))
	}

	return nil
}

// Ping pings server
func (c *Client) Ping() error {
	url := c.requestPath(PingCall)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("Ping: %s", err)
	}
	resp, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("Ping: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ping: %s", http.StatusText(resp.StatusCode))
	}

	return nil
}

// GetUsers gets all users from chat
func (c *Client) GetUsers() (users []User, err error) {
	url := c.requestPath(UsersCall)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("GetUsers: %s", err)
	}
	request.SetBasicAuth(c.user.AuthID, c.user.Token)

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("GetUsers: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetUsers: %s", http.StatusText(resp.StatusCode))
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("GetUsers: %s", err)
	}

	if err := json.Unmarshal(content, &users); err != nil {
		return nil, fmt.Errorf("GetUsers: %s", err)
	}

	return users, err
}

// GetRooms get all rooms from chat
func (c *Client) GetRooms() (rooms []Room, err error) {
	url := c.requestPath(RoomsCall)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("GetRooms: %s", err)
	}
	request.SetBasicAuth(c.user.AuthID, c.user.Token)

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("GetRooms: %s", err)
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("GetRooms: %s", err)
	}

	if err := json.Unmarshal(content, &rooms); err != nil {
		return nil, fmt.Errorf("GetRooms: %s", err)
	}

	return rooms, err
}

// RoomCreate creates new room and sets user as owner
func (c *Client) RoomCreate(name string) error {
	return nil
}

// RoomDelete deletes room from server
// note: user must be owner of a room
func (c *Client) RoomDelete(name string) error {
	return nil
}

// RoomJoin joins user to room
func (c *Client) RoomJoin(name string) error {
	return nil
}

// RoomExit exits user to room
// note: user can't be owner of the room
func (c *Client) RoomExit(name string) error {
	return nil
}
