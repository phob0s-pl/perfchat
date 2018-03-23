package chat

// Room represents place where users can chat
type Room struct {
	// Name is human readable name for room
	Name string
	// Creator is a name of user which created room
	Creator string
	// Users is list of users currently in chat
	Users []*User
}

// Join adds user to room
func (r *Room) Join(user *User) error {
	for _, roomUser := range r.Users {
		if user.Name == roomUser.Name {
			return ErrExists
		}
	}
	r.Users = append(r.Users, user)
	return nil
}

// Exit removes user from room
func (r *Room) Exit(user *User) error {
	for i, roomUser := range r.Users {
		if user.Name == roomUser.Name {
			r.Users = append(r.Users[:i], r.Users[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}
