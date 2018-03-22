package api

import "github.com/phob0s-pl/perfchat/chat"

const (
	// Version is API version
	Version = "1"
)

type API struct {
	engine chat.Chat
}


