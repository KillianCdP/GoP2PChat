package main

import (
	"GoP2PChat/config"
	"GoP2PChat/internals/host"
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"github.com/libp2p/go-libp2p-core/crypto"
	"os"
)

func main() {
	help := flag.Bool("help", false, "Display Help")
	cfg, err := config.ParseFlags()

	if err != nil {
		fmt.Printf("Could not read config")
		panic(err)
	}

	if *help {
		fmt.Printf("Simple example for peer discovery using mDNS. mDNS is great when you have multiple peers in local LAN.")
		fmt.Printf("Usage: \n   Run './chat-with-mdns'\nor Run './chat-with-mdns -host [host] -port [port] -rendezvous [string] -pid [proto ID]'\n")

		os.Exit(0)
	}

	fmt.Println("Started")

	ctx := context.Background()
	r := rand.Reader

	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		panic(err)
	}

	host.Start(ctx, prvKey, cfg)
}
