package main

import (
	"fmt"

	"github.com/pkg/errors"
)

func main() {
	fmt.Println(Boom())
}

func Boom() error {
	return errors.WithMessage(errors.New("boom"), "sorry")
}
