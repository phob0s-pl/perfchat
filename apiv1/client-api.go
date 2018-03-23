package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

// newRequest creates new requests and sets auth info
func (c *Client) newAPIRequest(method, apiCall string, body io.Reader) (*http.Request, error) {
	url := c.requestPath(apiCall)
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(c.user.AuthID, c.user.Token)
	return request, nil
}

// simpleDo makes http request and checks response code and returns body
func (c *Client) do(request *http.Request) ([]byte, error) {
	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", http.StatusText(resp.StatusCode))
	}

	return ioutil.ReadAll(resp.Body)
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

	request, err := c.newAPIRequest(http.MethodPost, UsersCall, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("AddUser: %s", err)
	}

	if _, err := c.do(request); err != nil {
		return fmt.Errorf("AddUser: %s", err)
	}
	return nil
}

// Ping pings server
func (c *Client) Ping() error {
	request, err := c.newAPIRequest(http.MethodGet, PingCall, nil)
	if err != nil {
		return fmt.Errorf("Ping: %s", err)
	}

	if _, err := c.do(request); err != nil {
		return fmt.Errorf("Ping: %s", err)
	}
	return nil
}

// GetUsers gets all users from chat
func (c *Client) GetUsers() (users []User, err error) {
	request, err := c.newAPIRequest(http.MethodGet, UsersCall, nil)
	if err != nil {
		return nil, fmt.Errorf("GetUsers: %s", err)
	}

	body, err := c.do(request)
	if err != nil {
		return users, fmt.Errorf("GetUsers: %s", err)
	}

	if err := json.Unmarshal(body, &users); err != nil {
		return nil, fmt.Errorf("GetUsers: %s", err)
	}

	return users, err
}

// GetRooms get all rooms from chat
func (c *Client) GetRooms() (rooms []Room, err error) {
	request, err := c.newAPIRequest(http.MethodGet, RoomsCall, nil)
	if err != nil {
		return nil, fmt.Errorf("GetRooms: %s", err)
	}

	body, err := c.do(request)
	if err != nil {
		return rooms, fmt.Errorf("GetRooms: %s", err)
	}

	if err := json.Unmarshal(body, &rooms); err != nil {
		return nil, fmt.Errorf("body: %s", err)
	}

	return rooms, err
}

// RoomCreate creates new room and sets user as owner
func (c *Client) RoomCreate(name string) error {
	payload, err := json.Marshal(&Room{
		Name: name,
	})
	if err != nil {
		return fmt.Errorf("RoomCreate: %s", err)
	}

	request, err := c.newAPIRequest(http.MethodPost, RoomsCreateCall, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("RoomCreate: %s", err)
	}

	if _, err := c.do(request); err != nil {
		return fmt.Errorf("RoomCreate: %s", err)
	}
	return nil
}

// RoomDelete deletes room from server
// note: user must be owner of a room
func (c *Client) RoomDelete(name string) error {
	payload, err := json.Marshal(&Room{
		Name: name,
	})
	if err != nil {
		return fmt.Errorf("RoomDelete: %s", err)
	}

	request, err := c.newAPIRequest(http.MethodPost, RoomsDeleteCall, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("RoomDelete: %s", err)
	}

	if _, err := c.do(request); err != nil {
		return fmt.Errorf("RoomDelete: %s", err)
	}
	return nil
}

// RoomJoin joins user to room
func (c *Client) RoomJoin(name string) error {
	payload, err := json.Marshal(&Room{
		Name: name,
	})
	if err != nil {
		return fmt.Errorf("RoomJoin: %s", err)
	}

	request, err := c.newAPIRequest(http.MethodPost, RoomsJoinCall, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("RoomJoin: %s", err)
	}

	if _, err := c.do(request); err != nil {
		return fmt.Errorf("RoomJoin: %s", err)
	}
	return nil
}

// RoomExit exits user to room
// note: user can't be owner of the room
func (c *Client) RoomExit(name string) error {
	payload, err := json.Marshal(&Room{
		Name: name,
	})
	if err != nil {
		return fmt.Errorf("RoomExit: %s", err)
	}

	request, err := c.newAPIRequest(http.MethodPost, RoomsExitCall, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("RoomExit: %s", err)
	}

	if _, err := c.do(request); err != nil {
		return fmt.Errorf("RoomExit: %s", err)
	}
	return nil
}
