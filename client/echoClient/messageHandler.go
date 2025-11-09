package main

import (
	"fmt"
	"log"
	"net"
)

func messageHandler(inputChannel <-chan string, doneChannel chan<- bool, conn net.Conn) {
	for event := range inputChannel {
		log.Printf("Client sending message: %s", event)
		_, err := fmt.Fprintln(conn, event)
		if err != nil {
			log.Printf("Error sending the message: %v", err)
			break
		}
		if event == "goodbye" {
			disconnected(conn)
			break
		}
	}

	doneChannel <- true
}
