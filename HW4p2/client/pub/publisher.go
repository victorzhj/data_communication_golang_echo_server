package main

func main() {
	inputChannel := make(chan string, 100)
	doneChannel := make(chan bool)

	go keyboardReader(inputChannel)

	go messageHandler(inputChannel, doneChannel)

	<-doneChannel
}
