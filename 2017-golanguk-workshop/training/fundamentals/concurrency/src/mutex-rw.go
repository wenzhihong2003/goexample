package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	quit := make(chan struct{})
	var mu sync.RWMutex
	var counter int

	set := func(i int) {
		mu.Lock()
		defer mu.Unlock()
		counter = i
	}

	get := func() int {
		mu.RLock()
		defer mu.RUnlock()
		return counter
	}

	go func() {
		tick := time.Tick(500 * time.Millisecond)
		for {
			select {
			case <-tick:
				i := get()
				fmt.Printf("counter: %d\n", i)
				if i == 10 {
					close(quit)
				}
			}
		}
	}()

	for i := 0; i <= 10; i++ {
		go set(i)
		time.Sleep(750 * time.Millisecond)
	}

	<-quit
}
