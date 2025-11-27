package main

import (
	"bufio"
	"fmt"
	"log"
	"strings"
	"time"
)

func messageHandler(inputChannel <-chan string, doneChannel chan<- bool) {
	conn := connectMain()
	defer disconnected(conn)

	buffer := []string{}
	messageId := 1

	for message := range inputChannel {
		messageSplit := strings.SplitN(message, " ", 2)
		if len(messageSplit) < 2 {
			log.Println("Please use format: <topic> <message>")
			continue
		}
		topic := messageSplit[0]
		payload := messageSplit[1]

		fullContent := fmt.Sprintf("%d:%s", messageId, payload)
		fullMessage := fmt.Sprintf("PUBLISH|%s|%s", topic, fullContent)

		if len(buffer) >= 5 {
			buffer = buffer[1:]
		}
		buffer = append(buffer, fullMessage)

		success := false
		for !success {
			_, err := fmt.Fprintln(conn, fullMessage)
			if err != nil {
				log.Println("Write error (Primary might be dead).")
			} else {
				conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))

				scanner := bufio.NewScanner(conn)
				if scanner.Scan() {
					conn.SetReadDeadline(time.Time{})
					success = true
				}
			}

			if !success {
				log.Println("Timeout or Error! Switching to Backup Broker...")
				conn.Close()

				conn = switchBackUp()

				log.Println("Resending last 5 messages to Backup...")
				for _, savedMsg := range buffer {
					fmt.Fprintln(conn, savedMsg)
					time.Sleep(50 * time.Millisecond)
				}
				success = true
			}
		}
		messageId++
	}

	doneChannel <- true
}
