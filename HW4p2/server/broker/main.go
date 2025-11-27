package main

import (
	"flag"
	"log"
	"net"
)

type ClientEvent struct {
	Packet ControlPacket
	Client net.Conn
}

func main() {
	mode := flag.String("mode", "primary", "Mode: 'primary' or 'backup'")
	port := flag.String("port", "8080", "Port to listen on")
	flag.Parse()

	log.Printf("Starting Broker in %s mode...", *mode)

	eventChannel := make(chan ClientEvent, 100)

	if *mode == "primary" {
		go messageHandler(eventChannel)
	} else {
		go backupHandler(eventChannel)
	}

	startConnectionHandler(eventChannel, *port)
}
