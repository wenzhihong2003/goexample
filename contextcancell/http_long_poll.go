package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	done := make(chan bool)
	closeNotifier, ok := w.(http.CloseNotifier)
	if !ok {
		panic("Expected http.ResponseWriter to be an http.CloseNotifier")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go doLongOperation(w, ctx, done)
	select {
	case <-done:
	case <-time.After(time.Second * 10):
		cancel()
		if !<-done {
			fmt.Fprint(w, "Server is busy.")
		}
	case <-closeNotifier.CloseNotify():
		cancel()
		fmt.Println("client has disconnected.")
		<-done
	}
}

func doLongOperation(w http.ResponseWriter, ctx context.Context, done chan bool) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		panic("Expected http.ResponseWriter to be an http.Flusher")
	}
	select {
	case <-time.After(time.Second * 3):
	case <-ctx.Done():
		done <- false
		return
	}
	fmt.Fprint(w, "Sup dude?")
	flusher.Flush()
	done <- true
}
