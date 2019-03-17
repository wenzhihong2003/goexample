package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

// 关闭多个goroutines, 使用context包
// 参考: https://chilts.org/2017/06/12/cancelling-multiple-goroutines

func main() {
	// create a context that we can cencel
	// 这个cancel函数可以多次调用, 没事
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// a waitgroup for the goroutines to tell us they've stopped
	wg := sync.WaitGroup{}

	// a channel for tick() to telll us they've stopped
	wg.Add(1)
	go tick(ctx, &wg)

	// a channel for tock() to tell us they've stopped
	wg.Add(1)
	go tock(ctx, &wg)

	// listen for ctrl-c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("main: received ctrl-c shutting down")

	// tell the goroutines to stop
	fmt.Println("main: telling goroutines to stop")
	cancel()

	// and wait for then both to reply back
	wg.Wait()
	fmt.Println("main: all goroutines have told us they've finished")
}

func tock(ctx context.Context, wg *sync.WaitGroup) {
	// tell the caller we've stopped
	defer wg.Done()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case now := <-ticker.C:
			fmt.Printf("tock: tock %s\n", now.UTC().Format("20060102-150405.000000000"))
		case <-ctx.Done():
			fmt.Println("tock: caller has told us to stop")
			return
		}
	}
}

func tick(ctx context.Context, wg *sync.WaitGroup) {
	// tell the caller we've stopped
	defer wg.Done()

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case now := <-ticker.C:
			fmt.Printf("tick: tick %s\n", now.UTC().Format("20060102-150405.000000000"))
		case <-ctx.Done():
			fmt.Println("tick: caller has told us to stop")
			return
		}
	}
}
