package config

import (
	"flag"
)

type Config struct {
	RendezvousString string
	ProtocolID       string
	ListenPort       int
}

func ParseFlags() (Config, error) {
	config := Config{}
	flag.IntVar(&config.ListenPort, "port", 3000, "Port to listen to")
	flag.StringVar(&config.RendezvousString, "rendezvous", "rendezvouspoint",
		"Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	flag.StringVar(&config.ProtocolID, "pid", "/chat/1.1.0", "Sets a protocol id for stream headers")
	flag.Parse()
	return config, nil
}
