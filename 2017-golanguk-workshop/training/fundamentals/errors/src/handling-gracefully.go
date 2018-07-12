package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	var f io.Reader
	var err error

	// try to read a file
	f, err = os.Open("/path/to/some/content.file")
	if err != nil {
		// create a fall back io.Reader so our program works
		f = bytes.NewBufferString("some fall back content")
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
