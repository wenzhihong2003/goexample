package main

import "sync"

func main() {
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i <5; i++ {
		go func(i int) {
			println(i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
