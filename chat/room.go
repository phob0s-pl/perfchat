package chat

// Room represents place where users can chat
type Room struct {
	// ID is unique identifier for room
	ID int
	// Name is human readable name for room
	Name string
	// Creator is a name of user which created room
	Creator string
	// Started contais date when room was created in format RFC3339
	Started string
}
