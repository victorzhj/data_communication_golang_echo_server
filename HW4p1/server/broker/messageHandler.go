package main

import (
	"fmt"
	"log"
	"net"
)

func messageHandler(events <-chan ClientEvent) {
	var subscribers map[string][]net.Conn = make(map[string][]net.Conn)
	log.Println("Broker is running...")
	for events := range events {
		cmd := events.Packet.Type
		topic := events.Packet.Topic
		payload := events.Packet.Payload
		client := events.Client
		switch cmd {
		case CMD_SUBSCRIBE:
			log.Printf("ACTION: Client %s subscribed to '%s'", client.RemoteAddr(), topic)
			subscribers[topic] = append(subscribers[topic], client)
		case CMD_PUBLISH:
			log.Printf("ACTION: Publishing to '%s': %s", topic, payload)
			targets, found := subscribers[topic]
			if !found || len(targets) == 0 {
				log.Println(" - No subscribers found for this topic.")
				continue
			}
			for _, sub := range targets {
				_, err := fmt.Fprintln(sub, payload)
				if err != nil {
					log.Printf("Error sending to subscriber %s: %v", sub.RemoteAddr(), err)
				}
			}
		default:
			log.Printf("Unknown command received: %s", cmd)
		}
	}
}
