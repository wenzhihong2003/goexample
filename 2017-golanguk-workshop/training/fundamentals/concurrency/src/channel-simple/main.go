package main

func main() {
	messages := make(chan string)

	// this go routine launches immediately
	go func() {
		// This line blocks until the channel is read from
		messages <- "hello!"
	}()

	// this line blocks until someone writes to the channel
	msg := <-messages

	println(msg)
}
