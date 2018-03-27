package main

import (
	"sync"
)

type Group struct {
	Creator string
	Name    string
	Users   []string
}

type Chat struct {
	sync.Mutex
	selfName string

	joinedGroups map[string]*Group
	knownGroups  map[string]*Group
	deletedRooms chan string
}

func NewChat(selfName string) *Chat {
	return &Chat{
		joinedGroups: make(map[string]*Group),
		knownGroups:  make(map[string]*Group),
		selfName:     selfName,
	}
}

func (c *Chat) GetRandomGroup() string {
	c.Lock()
	defer c.Unlock()
	if len(c.joinedGroups) == 0 {
		return ""
	}
	roomID := randSrc.Uint32() % uint32(len(c.joinedGroups))
	var i uint32
	for k := range c.joinedGroups {
		if i == roomID {
			return k
		}
		i++
	}
	return ""
}

func (c *Chat) ExitRandomGroup() string {
	c.Lock()
	defer c.Unlock()

	if len(c.joinedGroups) == 0 {
		return ""
	}

	roomID := randSrc.Uint32() % uint32(len(c.joinedGroups))
	var i uint32
	for k := range c.joinedGroups {
		if i == roomID {
			delete(c.joinedGroups, k)
			return k
		}
		i++
	}

	return ""
}

func (c *Chat) CreateRandomGroup() string {
	c.Lock()
	defer c.Unlock()

	name := RandString()
	c.joinedGroups[name] = &Group{
		Creator: c.selfName,
		Name:    name,
		Users:   []string{c.selfName},
	}
	return name
}

func (c *Chat) JoinRandomGroup() string {
	c.Lock()
	defer c.Unlock()

	if len(c.knownGroups) == 0 {
		return ""
	}
	roomID := randSrc.Uint32() % uint32(len(c.knownGroups))
	var i uint32
	for key, group := range c.knownGroups {
		if i == roomID {
			c.joinedGroups[key] = group
			delete(c.knownGroups, key)
			return key
		}
		i++
	}
	return ""
}

func (c *Chat) IsGroupMember(userName, roomName string) bool {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.joinedGroups[roomName]; ok {
		group := c.joinedGroups[roomName]
		for _, user := range group.Users {
			if user == userName {
				return true
			}
		}
	}
	return false
}

func (c *Chat) AddKnownGroup(name, creator string) {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.knownGroups[name]; !ok {
		c.knownGroups[name] = &Group{
			Creator: creator,
			Name:    name,
			Users:   []string{creator},
		}
	}
}

func (c *Chat) AddUserToGroup(groupname, username string) {
	c.Lock()
	defer c.Unlock()

	if group, ok := c.joinedGroups[groupname]; ok {
		for _, user := range group.Users {
			if user == username {
				return
			}
		}
		group.Users = append(group.Users, username)
		return
	}

	if group, ok := c.knownGroups[groupname]; ok {
		for _, user := range group.Users {
			if user == username {
				return
			}
		}
		group.Users = append(group.Users, username)
		return
	}
}

func (c *Chat) ExitGroup(roomname, username string) {
	c.Lock()
	defer c.Unlock()

	if group, ok := c.joinedGroups[roomname]; ok {
		if group.Creator == username {
			delete(c.joinedGroups, roomname)
			return
		}
		for i, user := range group.Users {
			if user == username {
				group.Users = append(group.Users[:i], group.Users[i+1:]...)
				return
			}
		}
	}

	if group, ok := c.knownGroups[roomname]; ok {
		if group.Creator == username {
			delete(c.joinedGroups, roomname)
			return
		}
		for i, user := range group.Users {
			if user == username {
				group.Users = append(group.Users[:i], group.Users[i+1:]...)
				return
			}
		}
	}
}
