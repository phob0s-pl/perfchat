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

func (a *API) CreateRoomRoute() *Route {
	return &Route{
		HandlerFunc: a.CreateRoom,
		Method:      http.MethodPost,
		Name:        "CreateRoom",
		Pattern:     GetPath(RoomsCreateCall),
		WithPrefix:  false,
	}
}

func (a *API) JoinRoomRoute() *Route {
	return &Route{
		HandlerFunc: a.RoomJoin,
		Method:      http.MethodPost,
		Name:        "JoinRoom",
		Pattern:     GetPath(RoomsJoinCall),
		WithPrefix:  false,
	}
}

func (a *API) ExitRoomRoute() *Route {
	return &Route{
		HandlerFunc: a.RoomExit,
		Method:      http.MethodPost,
		Name:        "ExitRoom",
		Pattern:     GetPath(RoomsExitCall),
		WithPrefix:  false,
	}
}

func (a *API) SendMessageRoute() *Route {
	return &Route{
		HandlerFunc: a.SendMessage,
		Method:      http.MethodGet,
		Name:        "SendMessage",
		Pattern:     GetPath(MessageCall),
		WithPrefix:  false,
	}
}

func (a *API) ReceiveMessageRoute() *Route {
	return &Route{
		HandlerFunc: a.ReceiveMessage,
		Method:      http.MethodPost,
		Name:        "ReceiveMessage",
		Pattern:     GetPath(MessageCall),
		WithPrefix:  false,
	}
}
