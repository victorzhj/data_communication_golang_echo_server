package main

import (
	"log"
	"net"
)

const brokerMainAddr = "127.0.0.1:8080"
const brokerBackUpAddr = "127.0.0.1:8081"

func connectMain() net.Conn {
	conn := connect(brokerMainAddr)
	return conn
}

func switchBackUp() (conn net.Conn) {
	conn = connect(brokerBackUpAddr)
	return conn
}

func connect(address string) net.Conn {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("Couldn't connect to the server: %v", err)
	}
	return conn
}

func disconnected(conn net.Conn) {
	conn.Close()
	log.Printf("Disconnected from server: %v", conn.RemoteAddr())
}
