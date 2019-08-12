package host

import (
	"bufio"
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
	"os"

	"GoP2PChat/config"
	"github.com/multiformats/go-multiaddr"
)

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

	peerChan := initMDNS(ctx, host, config.RendezvousString)

	peer := <-peerChan // Block until we discover a peer
	fmt.Println("Found peer:", peer, ", connecting")

	if err := host.Connect(ctx, peer); err != nil {
		fmt.Println("Connection failed:", err)
	}

	stream, err := host.NewStream(ctx, peer.ID, protocol.ID(config.ProtocolID))

	if err != nil {
		fmt.Println("Stream open failed:", err)
	} else {
		rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

		go writeData(rw)
		go readData(rw)
		fmt.Println("Connected to :", peer)
	}

	select {}
}

func handleStream(stream network.Stream) {
	fmt.Println("Got a new stream !")

	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go readData(rw)
	go writeData(rw)
}

func readData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from buffer")
			panic(err)
		}

		if str == "" {
			return
		}

		if str != "\n" {
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
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
