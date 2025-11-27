package main

import (
	"bufio"
	"log"
	"net"
)

func startConnectionHandler(events chan<- ClientEvent, port string) {
	address := ":" + port
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to start listener on %s: %v", address, err)
	}
	defer listener.Close()

	log.Printf("Listening on %s...", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept: %v", err)
			continue
		}
		go handleClientConnection(conn, events)
	}
}

func handleClientConnection(conn net.Conn, events chan<- ClientEvent) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		packet, err := parsePacket(message)
		if err != nil {
			log.Printf("Bad packet from %s: %v", conn.RemoteAddr(), err)
			continue
		}
		events <- ClientEvent{Packet: packet, Client: conn}
	}
}
