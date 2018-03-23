package api

import "net/http"

func (a *API) AddUserRoute() *Route {
	return &Route{
		HandlerFunc: a.AddUser,
		Method:      http.MethodPost,
		Name:        "AddUser",
		Pattern:     GetPath(UsersCall),
		WithPrefix:  true,
	}
}

func (a *API) PingRoute() *Route {
	return &Route{
		HandlerFunc: a.Ping,
		Method:      http.MethodGet,
		Name:        "Ping",
		Pattern:     GetPath(PingCall),
		WithPrefix:  false,
	}
}

func (a *API) GetUsersRoute() *Route {
	return &Route{
		HandlerFunc: a.GetUsers,
		Method:      http.MethodGet,
		Name:        "GetUsers",
		Pattern:     GetPath(UsersCall),
		WithPrefix:  false,
	}
}

func (a *API) GetRoomsRoute() *Route {
	return &Route{
		HandlerFunc: a.GetRooms,
		Method:      http.MethodGet,
		Name:        "GetRooms",
		Pattern:     GetPath(RoomsCall),
		WithPrefix:  false,
	}
}
