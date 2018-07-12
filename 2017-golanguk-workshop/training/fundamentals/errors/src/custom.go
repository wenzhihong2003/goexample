package main

import (
	"fmt"
)

type Yoko struct{}

func (Yoko) Error() string {
	return "i broke up the Beatles"
}

func main() {
	for _, b := range []string{"Paul", "George", "John", "Ringo"} {
		err := play(b)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
	// Paul
	// George
	// i broke up the Beatles
}

func play(b string) error {
	if b == "John" {
		return Yoko{}
	}
	fmt.Println(b)
	return nil
}
