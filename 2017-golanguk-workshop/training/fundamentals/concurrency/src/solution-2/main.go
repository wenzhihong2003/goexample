package main

func process(messages chan string, quit chan struct{}) {
	for m := range messages {
		if m == "" {
			break
		}
		println(m)
	}

	close(quit)
}

func main() {
	messages := make(chan string, 5)
	quit := make(chan struct{})

	go process(messages, quit)

	fruits := []string{"apple", "plum", "peach", "pear", "grape", ""}
	for _, s := range fruits {
		messages <- s
	}

	<-quit
	println("done")
}
