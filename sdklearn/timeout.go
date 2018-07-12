package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// https://gocn.vip/question/1979
// 如何在父 goroutine 中通过超时控制来结束子 goroutine

func main() {
	method2()
}

func method2() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			fmt.Println(ctx.Err())
			return
		case <-time.NewTicker(10 * time.Second).C:
			fmt.Println("Done after 10 seconds")
		}
	}()

	wg.Wait()

}

func method1() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		ch := make(chan bool)
		go func() {
			time.Sleep(10 * time.Second)
			ch <- true
		}()

		select {
		case <-ctx.Done():
			fmt.Println(ctx.Err())
		case <-ch:
			fmt.Println("Done after 10 seconds")
		}
	}(ctx)
	wg.Wait()
}
