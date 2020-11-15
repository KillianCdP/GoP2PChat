package user

import (
	"GoP2PChat/internals/message"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
)

type User struct {
	AddressInfo peer.AddrInfo
	Name string
	stream network.Stream
}

func New (addressInfo peer.AddrInfo, name string, stream network.Stream) User {
	return User{
		AddressInfo: addressInfo,
		Name:        name,
		stream:      stream,
	}
}

func (u User) SendMessage(message message.Message) {
	w := bufio.NewWriter(u.stream)
	bytes, err := json.Marshal(message)
	if err != nil {
		fmt.Print("Could not unmarshal message ", message)
	}
	_, err = w.Write(bytes)
	if err != nil {
		fmt.Println("Could not write message ", message)
	}

	err = w.Flush()
	if err != nil {
		fmt.Println("Error flushing data")
	}
}


func GenerateUserName() string {
	return namesgenerator.GetRandomName(0)
}