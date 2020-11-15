package host

import (
	"GoP2PChat/internals/message"
	"GoP2PChat/internals/user"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"os"
	"strings"

	"GoP2PChat/config"
	"github.com/multiformats/go-multiaddr"
)

var userName string

var users []user.User

// Start starts the client
func Start(ctx context.Context, prvKey crypto.PrivKey, config config.Config) {

	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", "0.0.0.0", config.ListenPort))

	host, err := libp2p.New(
		ctx,
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)

	if err != nil {
		panic(err)
	}

	host.SetStreamHandler(protocol.ID(config.ProtocolID), handleStream)

	fmt.Printf("\n[*] Your MultiAddress is: /ip4/%s/tcp/%d/p2p/%s\n", "0.0.0.0", config.ListenPort, host.ID().Pretty())

	userName = user.GenerateUserName()
	fmt.Printf("\n[*] Your user name is %s\n ", userName)

	peerChan := initMDNS(ctx, host, config.RendezvousString)

	for addrInfo := range peerChan {
		fmt.Println("Found addrInfo:", addrInfo, ", connecting")

		if err := host.Connect(ctx, addrInfo); err != nil {
			fmt.Println("Connection failed:", err)
		}

		stream, err := host.NewStream(ctx, addrInfo.ID, protocol.ID(config.ProtocolID))

		if err != nil {
			fmt.Println("Stream open failed:", err)
		} else {
			addUser(addrInfo, stream)
			handleStream(stream)
		}
	}

	listenInput()

	select {}
}

func handleStream(stream network.Stream) {
	fmt.Println("New stream !")
	r := bufio.NewReader(bufio.NewReader(stream))
	go readData(r)
}

func addUser(aI peer.AddrInfo, s network.Stream) {
	users = append(users, user.New(aI, "Anonymous", s))
}

func readData(r *bufio.Reader) {
	messageUnmarshaller := json.NewDecoder(r)
	message := message.Message{}
	for {
		//str, err := r.ReadString('\n')
		err := messageUnmarshaller.Decode(&message)
		if err != nil {
			fmt.Println("Could not unmarshal received message")
		}

		fmt.Printf("Received message %+v\n", message)
	}
}

func listenInput() {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin")
			panic(err)
		}

		parseInput(sendData)
	}
}

func parseInput(input string) {
	if !strings.HasPrefix(input, "/") {
		return
	}
	input = strings.TrimSuffix(input, "/")
	params := strings.Split(input, " ")

	if len(params) < 2 {
		fmt.Println("Command too short")
		return
	}

	switch params[0] {
	case "message":
		u, err := getUser(users, params[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		u.SendMessage(message.Message{
			MessageType: message.TextMessage,
			Message:     strings.Join(params[2:], " "),
			RoomName:    "",
		})
	case "broadcast":
		m := message.Message{
			MessageType: message.Broadcast,
			Message:     strings.Join(params[2:], " "),
			RoomName:    "",
		}
		for _, u := range(users) {
			u.SendMessage(m)
		}

	}
}

func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin")
			panic(err)
		}

		_, err = rw.WriteString(fmt.Sprintf("%s\n", sendData))
		if err != nil {
			fmt.Println("Error writing to buffer")
			panic(err)
		}

		err = rw.Flush()
		if err != nil {
			fmt.Println("Error flushing buffer")
			panic(err)
		}
	}
}