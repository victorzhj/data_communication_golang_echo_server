package main

func main() {
	inputChannel := make(chan string)
	doneChannel := make(chan bool)

	conn := connect()

	go keyboardReader(inputChannel)

	go messageHandler(inputChannel, doneChannel, conn)

	go dispServerReply(conn)

	<-doneChannel
}
