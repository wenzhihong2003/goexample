package main

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const N = 50

var found = make(chan string)

func main() {
	for i := 0; i < N; i++ {
		go search(i)
	}
	fmt.Print(<-found)
	time.Sleep(3 * time.Second)
}

func search(i int) {
	for {
		fmt.Print(".")
		x := rand.Intn(N)
		if x == 42 {
			fmt.Print("found")
			found <- fmt.Sprintf("routine %d: found 42!", i)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}
