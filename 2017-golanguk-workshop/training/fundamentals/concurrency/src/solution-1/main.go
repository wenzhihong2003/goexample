package main

import (
	"fmt"
	"sync"
)

func count(prefix string, count int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < count; i++ {
		fmt.Println(prefix, i)
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go count("first: ", 50, &wg)
	go count("second: ", 50, &wg)
	wg.Wait()
}
