package main

import "fmt"

type Command struct {
	ID     int
	Result string
}

func (c Command) Error() string {
	return fmt.Sprintf("%s %d", c.Result, c.ID)
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
