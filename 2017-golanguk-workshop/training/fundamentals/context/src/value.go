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
	ctx = context.WithValue(ctx, "magic", 42)
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
	magic := ctx.Value("magic").(int)
Loop:
	for {
		select {
		case <-ctx.Done():
			fmt.Print("x")
			break Loop
		default:
			fmt.Print(".")
			x := rand.Intn(N)
			if x == magic {
				fmt.Print("found")
				found <- fmt.Sprintf("routine %d: found %d!", i, magic)
				break Loop
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}
