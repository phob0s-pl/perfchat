package chat

const (
	// AdminRole represents administrative role
	AdminRole = "admin"
	// UserRole represents regular user role
	UserRole = "user"
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
	// AuthID is username for authorization
	AuthID string
	// Token is password for authorization
	Token string
}

// CanAddUser checks whether user can add another one
func (u *User) CanAddUser() bool {
	return u.Role == AdminRole
}
