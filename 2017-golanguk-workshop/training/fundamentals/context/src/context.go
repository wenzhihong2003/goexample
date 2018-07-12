package main

import (
	"context"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for i := 0; i < N; i++ {
		go search(ctx, i)
	}
	m := <-found
	cancel()
	fmt.Print(m)
	time.Sleep(3 * time.Second)
}

func search(ctx context.Context, i int) {
Loop:
	for {
		select {
		case <-ctx.Done():
			fmt.Print("x")
			break Loop
		default:
			fmt.Print(".")
			x := rand.Intn(N)
			if x == 42 {
				fmt.Print("found")
				found <- fmt.Sprintf("routine %d: found 42!", i)
				break Loop
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}
