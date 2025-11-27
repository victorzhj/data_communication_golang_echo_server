package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

// Config
const (
	primaryAddr = "127.0.0.1:8080"
	backupAddr  = "127.0.0.1:8081"
)

func main() {
	topic := "topicA" // Hardcoded for simplicity, or use os.Args[1]

	// 1. Connect to Primary initially
	currentBroker := primaryAddr
	conn, err := net.Dial("tcp", currentBroker)
	if err != nil {
		log.Fatalf("Could not connect to primary: %v", err)
	}
	defer conn.Close()

	// 2. Setup State
	seq := 1
	historyBuffer := []string{} // Stores the payloads (e.g. "1", "2")

	// Loop forever (10 Hz = 10 times per second)
	for {
		// A. Prepare the Message
		payload := fmt.Sprintf("%d", seq)

		// B. Update History Buffer
		// Add new message to the end
		historyBuffer = append(historyBuffer, payload)
		// If we have more than 5, remove the oldest (from the front)
		if len(historyBuffer) > 5 {
			historyBuffer = historyBuffer[1:]
		}

		// C. Send and Wait for ACK
		success := sendAndCheckAck(conn, topic, payload)

		if !success {
			// CRASH DETECTED!
			log.Printf("Primary %s failed! Switching to Backup...", currentBroker)
			conn.Close()

			// D. Switch to Backup
			currentBroker = backupAddr
			conn, err = net.Dial("tcp", currentBroker)
			if err != nil {
				log.Fatalf("Backup also failed! Exiting: %v", err)
			}

			// E. RESEND HISTORY (The critical requirement)
			log.Println("Resending last 5 messages...")
			for _, oldPayload := range historyBuffer {
				// We don't necessarily wait for ACKs during recovery resend
				// strictly based on the prompt, but it's safer to just send them.
				packet := fmt.Sprintf("PUBLISH|%s|%s", topic, oldPayload)
				fmt.Fprintln(conn, packet)
				// Small delay to ensure they don't stick together
				time.Sleep(10 * time.Millisecond)
			}

			// We successfully switched. The next loop will send normal messages to Backup.
		}

		// F. Prepare for next message
		seq++
		time.Sleep(100 * time.Millisecond) // 10Hz rate
	}
}

// Reuse the helper function we discussed earlier
func sendAndCheckAck(conn net.Conn, topic string, payload string) bool {
	// Construct packet
	msg := fmt.Sprintf("PUBLISH|%s|%s", topic, payload)

	// Send
	_, err := fmt.Fprintln(conn, msg)
	if err != nil {
		return false
	}

	// Set Deadline (500ms)
	conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))

	// Wait for ACK
	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		// We got data back (assuming it's "ACK")
		conn.SetReadDeadline(time.Time{}) // Clear deadline
		return true
	}

	// Timeout or Error
	return false
}
