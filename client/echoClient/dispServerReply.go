package main

import (
	"bufio"
	"fmt"
	"net"
)

func dispServerReply(conn net.Conn) {
	reader := bufio.NewScanner(conn)
	for reader.Scan() {
		fmt.Println("Echo: ", reader.Text())
	}
}
