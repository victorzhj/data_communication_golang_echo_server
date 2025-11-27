package main

import (
	"fmt"
	"log"
)

func messageHandler(events <-chan ClientEvent) {
	log.Println("Message handler (brain) is running.")

	for event := range events {

		log.Printf("Brain processing '%s' from %s", event.Message, event.Client.RemoteAddr())

		if event.Message == "goodbye" {
			log.Printf("Brain: Client %s said goodbye. Closing connection.", event.Client.RemoteAddr())
			fmt.Fprintln(event.Client, "Goodbye!")
			event.Client.Close()
		} else {
			// Echo the message back to the client because in go
			// connections are treated as files. Literally writing
			//  the echo back to the file (client)
			_, err := fmt.Fprintln(event.Client, event.Message)
			if err != nil {
				log.Printf("Error writing to client %s: %v", event.Client.RemoteAddr(), err)
			}
		}
	}
}
