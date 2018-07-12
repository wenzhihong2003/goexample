package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// https://gist.github.com/cee-dub/883789dc11c82ae5d05f

type SSE struct{}

func (*SSE) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "cannot stream", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("connection", "keep-alive")

	cn, ok := w.(http.CloseNotifier)
	if !ok {
		http.Error(w, "cannot stream", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case <-cn.CloseNotify():
			log.Println("done: closed connection")
		case msg := <-messages:
			fmt.Fprintf(w, "data:%s\n\n", msg)
			f.Flush()
		}
	}
}

var messages = make(chan string)

func main() {
	http.Handle("/", &SSE{})
	go func() {
		for i := 0; i < 10; i++ {
			messages <- "yo"
			time.Sleep(time.Second)
		}
	}()
	log.Fatal(http.ListenAndServe(":3000", nil))
}
