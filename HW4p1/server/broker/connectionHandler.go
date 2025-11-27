package main

import (
	"bufio"
	"log"
	"net"
)

func startConnectionHandler(events chan<- ClientEvent) {
	// Listen for incoming TCP connections on port 8080.
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to start listener: %v", err)
	}
	defer listener.Close()

	log.Println("Connection handler is listening...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		log.Printf("Accepted new client: %s", conn.RemoteAddr())
		// Start a new goroutine to handle this client.
		go handleClientConnection(conn, events)
	}
}

func handleClientConnection(conn net.Conn, events chan<- ClientEvent) {
	// Ensure the connection is closed when this goroutine exits.
	defer conn.Close()

	// Create a scanner to read messages line by line.
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		message := scanner.Text()
		packet, err := parsePacket(message)
		if err != nil {
			log.Printf("Parse error from %s: %v", conn.RemoteAddr(), err)
			continue
		}
		log.Printf("Proxy read '%s' from %s", message, conn.RemoteAddr())

		event := ClientEvent{
			Packet: packet,
			Client: conn,
		}
		// Send the read message to the message handler.
		events <- event
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from client %s: %v", conn.RemoteAddr(), err)
	} else {
		log.Printf("Client %s disconnected.", conn.RemoteAddr())
	}
}
