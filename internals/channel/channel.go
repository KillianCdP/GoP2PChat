package channel

import (
	"GoP2PChat/internals/message"
	"GoP2PChat/internals/user"
)

type Channel struct {
	Name string
	Members []user.User
}

func SendMessage(message message.Message, channel Channel) {
	for _, u := range channel.Members {
		u.SendMessage(message)
	}
}