package chat

type Chat struct {
	users []*User
}

func (c *Chat) AddUser(user *User) {
	c.users = append(c.users, user)
}
