package main

import (
	"log"
	"net"
)

// Struct to represent an message from a client.
type ClientEvent struct {
	Packet ControlPacket
	Client net.Conn
}

func main() {
	log.Println("Starting reactor server on :8080...")

	// Create a channel for messages to pass to the message handler.
	eventChannel := make(chan ClientEvent)

	// Start the message handler goroutine
	go messageHandler(eventChannel)

	startConnectionHandler(eventChannel)
}
