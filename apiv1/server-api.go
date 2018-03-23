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

func (a *API) GetUsers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	_, err := a.getRequestingUser(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
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

func (a *API) AddUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		user = &User{}
	)

	requester, err := a.getRequestingUser(r)
	if err != nil || !requester.CanAddUser() {
		w.WriteHeader(http.StatusUnauthorized)
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

func (a *API) Ping(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
}
