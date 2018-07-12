package main

import "sync"

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	var mu sync.Mutex

	var count int

	go func() {
		for i := 0; i < 5; i++ {
			mu.Lock()
			count = i
			println(i)
			mu.Unlock()
		}
		wg.Done()
	}()

	wg.Wait()
}
