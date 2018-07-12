package main

import "fmt"

// Task: Implement the error interface on the Command type

type Command struct {
	ID     int
	Result string
}

func main() {
	fmt.Println("Starting")
	if err := process(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Completed")
}

func process() error {
	c := Command{ID: 1, Result: "unable to initialize command"}
	return c
}
