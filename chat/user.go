package chat

const (
	adminRole = "admin"
	userRole  = "user"
)

// User describes single user in chat
type User struct {
	// Name is unique user name
	Name string
	// Role determines user permissions
	// Currently supported are:
	// - admin : can add users +
	// - user - can chat, create and delete rooms
	Role string
}

// CanAddUser checks whether user can add another one
func (u *User) CanAddUser() bool {
	return u.Role == adminRole
}
