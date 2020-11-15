package message

// Message types
const (
	TextMessage = "TextMessage"
	Broadcast = "Broadcast"
	IdentityRequest = "IdentityRequest"
	IdentityReply = "IdentityReply"
)

type Message struct {
	MessageType string
	Message     string
	RoomName    string
}