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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	for i := 0; i < N; i++ {
		go search(ctx, i)
	}
	select {
	case m := <-found:
		cancel()
		fmt.Print(m)
	case <-ctx.Done():
		cancel()
		if err := ctx.Err(); err != nil {
			fmt.Println(err)
		}
	}
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
			// number changed to -1 to force a timeout.
			if x == -1 {
				fmt.Print("found")
				found <- fmt.Sprintf("routine %d: found 42!", i)
				break Loop
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}
