package main

import (
	"errors"
	"log"
)

func main() {
	if err := boom(); err != nil {
		log.Fatal(err) // boom!
	}
}

func boom() error {
	return errors.New("boom!")
}

func greetOrBoom() (string, error) {
	return "hello", errors.New("boom!")
}
