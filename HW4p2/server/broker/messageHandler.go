package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

const backupBrokerAddr = "127.0.0.1:8081"

func messageHandler(events <-chan ClientEvent) {
	subscribers := make(map[string][]net.Conn)

	var backupConn net.Conn
	var err error

	log.Println("Primary Broker is running...")

	for event := range events {
		cmd := event.Packet.Type
		topic := event.Packet.Topic
		payload := event.Packet.Payload
		client := event.Client

		switch cmd {
		case CMD_SUBSCRIBE:
			log.Printf("ACTION: Client %s subscribed to '%s'", client.RemoteAddr(), topic)
			subscribers[topic] = append(subscribers[topic], client)

		case CMD_PUBLISH:
			log.Printf("ACTION: Processing PUBLISH to '%s'", topic)

			if backupConn == nil {
				backupConn, err = net.Dial("tcp", backupBrokerAddr)
				if err != nil {
					log.Printf("Warning: Backup Broker offline: %v", err)
				}
			}

			if backupConn != nil {
				msg := fmt.Sprintf("%s|%s|%s\n", CMD_BACKUP, topic, payload)
				_, err := fmt.Fprintf(backupConn, msg)
				if err != nil {
					log.Printf("Failed to replicate: %v", err)
					backupConn.Close()
					backupConn = nil
				}
			}

			delay := time.Duration(rand.Intn(100)+50) * time.Millisecond
			time.Sleep(delay)

			targets, found := subscribers[topic]
			if found && len(targets) > 0 {
				for _, sub := range targets {
					_, err := fmt.Fprintln(sub, payload)
					if err != nil {
						log.Printf("Error sending to sub: %v", err)
					}
				}
			}

			if backupConn != nil {
				msg := fmt.Sprintf("%s|%s|%s\n", CMD_CLEAR, topic, payload)
				fmt.Fprintf(backupConn, msg)
			}

			fmt.Fprintln(client, "ACK")

		default:
			log.Printf("Unknown command received: %s", cmd)
		}
	}
}
