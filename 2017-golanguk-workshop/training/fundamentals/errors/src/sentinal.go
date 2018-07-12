package main

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
)

var ErrOver = errors.New("OVER")

func main() {
	for {
		b, err := Next()
		if err != nil {
			if errors.Cause(err) != ErrOver {
				log.Fatal(err)
			}
			break
		}
		fmt.Println(b)
	}
}

var beatles = []string{"John", "Paul", "George", "Ringo"}
var index int

func Next() (string, error) {
	defer func() { index++ }()
	if index >= len(beatles) {
		return "", errors.Wrap(ErrOver, "ran out of beatles")
	}

	return beatles[index], nil
}
