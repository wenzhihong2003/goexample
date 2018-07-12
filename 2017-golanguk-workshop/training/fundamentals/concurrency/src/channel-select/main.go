package main

import "os"

func main() {
	c := make(chan string)
	quit := make(chan struct{})

	go func(messages []string) {
		for _, s := range messages {
			c <- s
		}
		close(quit)
	}([]string{"hi", "bye"})

	for {
		select {
		case message := <-c:
			println(message)
		case <-quit:
			println("shutting down")
			os.Exit(0)
		}
	}
}
