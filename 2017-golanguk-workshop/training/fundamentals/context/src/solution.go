package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

const hamlet = "http://hamlet.gopherguides.com"

type Result struct {
	LineNumber int
	Text       string
}

func (r Result) String() string {
	return fmt.Sprintf("%d: %s", r.LineNumber, r.Text)
}

func search(ctx context.Context, book []byte, results chan Result) {
Loop:
	for i, line := range bytes.Split(book, []byte("\n")) {
		select {
		case <-ctx.Done():
			break Loop
		default:
			if bytes.Contains(line, []byte("Hamlet")) {
				results <- Result{
					LineNumber: i + 1,
					Text:       string(line),
				}
			}
		}
	}
}

func retrieveBook() ([]byte, error) {
	res, err := http.Get(hamlet)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if res.StatusCode != 200 {
		return nil, errors.Errorf("received a non-success message %d", res.StatusCode)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return b, errors.WithStack(err)
	}

	return b, nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	book, err := retrieveBook()
	if err != nil {
		log.Fatal(err)
	}

	results := make(chan Result)

	go search(ctx, book, results)

Loop:
	for i := 0; i < 50; i++ {
		select {
		case r := <-results:
			fmt.Printf("[%d] %s\n", i, r)
		case <-ctx.Done():
			cancel()
			if err := ctx.Err(); err != nil {
				log.Fatal(err)
			}
			break Loop
		}
	}
}
