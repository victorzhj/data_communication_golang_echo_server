package main

import (
	"bufio"
	"fmt"
	"net"
)

func acknowledger(conn net.Conn) {
	reader := bufio.NewScanner(conn)
	for reader.Scan() {
		fmt.Println("Echo: ", reader.Text())
	}
}
