package main

import (
	"bufio"
	"fmt"
	"os"
)

func keyboardReader(inputChannel chan<- string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Usage: type <topic> <message>:")
	for scanner.Scan() {
		text := scanner.Text()
		inputChannel <- text
	}
}
