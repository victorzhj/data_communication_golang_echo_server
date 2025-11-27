package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const brokerAddr = "127.0.0.1:8080"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run subscriber.go <topic>")
		return
	}
	topic := os.Args[1]

	conn, err := net.Dial("tcp", brokerAddr)
	if err != nil {
		log.Fatalf("Failed to connect to broker: %v", err)
	}
	defer conn.Close()

	subMsg := fmt.Sprintf("SUBSCRIBE|%s|", topic)

	_, err = fmt.Fprintln(conn, subMsg)
	if err != nil {
		log.Fatalf("Failed to send subscription: %v", err)
	}
	fmt.Printf("Subscribed to '%s'. Waiting for messages...\n", topic)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()

		fmt.Printf("Received: %s\n", msg)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Connection error: %v", err)
	}
}
