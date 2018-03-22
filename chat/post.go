package chat

type Post struct {
	// Created is timestamp of post creation
	Created int64
	// User is name of user who created post
	User string
	// Room is name of room where post is designated
	Room string
	// Message is actual message content
	Message string
}
