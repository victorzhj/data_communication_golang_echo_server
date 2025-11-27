package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const brokerAddr = "127.0.0.1:8080"

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run publisher.go <topic> <message>")
		return
	}
	topic := os.Args[1]
	payload := strings.Join(os.Args[2:], " ")

	conn, err := net.Dial("tcp", brokerAddr)
	if err != nil {
		log.Fatalf("Failed to connect to broker: %v", err)
	}
	defer conn.Close()

	msg := fmt.Sprintf("PUBLISH|%s|%s", topic, payload)

	_, err = fmt.Fprintln(conn, msg)
	if err != nil {
		log.Fatalf("Error publishing message: %v", err)
	}

	log.Printf("Published to '%s': %s", topic, payload)
}
