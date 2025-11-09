package main

import (
	"log"
	"net"
)

func connect() net.Conn {
	conn, err := net.Dial("tcp", "192.168.1.140:8080")
	if err != nil {
		log.Fatalf("Couldn't connect to the server: %v", err)
	}
	return conn
}

func disconnected(conn net.Conn) {
	conn.Close()
	log.Printf("Disconnected from server: %v", conn.RemoteAddr())
}
