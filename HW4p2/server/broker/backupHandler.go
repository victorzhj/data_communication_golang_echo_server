package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

const EVENT_PRIMARY_DEAD = "INTERNAL_PRIMARY_DEAD"

func backupHandler(events chan ClientEvent) {
	subscribers := make(map[string][]net.Conn)
	messageStore := make(map[string][]string)

	isPrimaryAlive := true

	// WATCHDOG
	go func() {
		for {
			time.Sleep(1 * time.Second)

			conn, err := net.Dial("tcp", "127.0.0.1:8080")
			if err != nil {
				log.Println("!!! PRIMARY BROKER DETECTED DEAD !!!")
				dummyPacket := ControlPacket{Type: EVENT_PRIMARY_DEAD}
				events <- ClientEvent{Packet: dummyPacket}
				return
			}
			conn.Close()
		}
	}()

	log.Println("Backup Handler started. Monitoring Primary...")

	for event := range events {
		cmd := event.Packet.Type
		topic := event.Packet.Topic
		payload := event.Packet.Payload
		client := event.Client

		switch cmd {
		case CMD_SUBSCRIBE:
			subscribers[topic] = append(subscribers[topic], client)
			log.Printf("Backup: Client subscribed to %s", topic)

		case CMD_BACKUP:
			messageStore[topic] = append(messageStore[topic], payload)

		case CMD_CLEAR:
			stored := messageStore[topic]
			if len(stored) > 0 {
				messageStore[topic] = stored[1:]
			}

		case EVENT_PRIMARY_DEAD:
			isPrimaryAlive = false
			log.Println("Backup: Taking over as Primary! Flushing stored messages...")

			for t, msgs := range messageStore {
				for _, msg := range msgs {
					for _, sub := range subscribers[t] {
						fmt.Fprintln(sub, msg)
					}
				}
				messageStore[t] = []string{}
			}

		case CMD_PUBLISH:
			if !isPrimaryAlive {
				log.Printf("Backup (Active): Publishing to %s: %s", topic, payload)

				time.Sleep(100 * time.Millisecond)

				for _, sub := range subscribers[topic] {
					fmt.Fprintln(sub, payload)
				}

				fmt.Fprintln(client, "ACK")
			}
		}
	}
}
