package main

import (
	"bufio"
	"fmt"
	"os"
)

func keyboardReader(inputChannel chan<- string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Type a message (or 'goodbye' to quit):")
	for scanner.Scan() {
		text := scanner.Text()
		inputChannel <- text
	}
}
