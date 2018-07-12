package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	var urls = []string{
		"http://www.gopherguides.com/",
		"http://www.golang.org/",
		"http://www.google.com/",
	}
	// Increment the WaitGroup counter.
	wg.Add(len(urls))
	for i, url := range urls {
		// Launch a goroutine to fetch the URL.
		go func(i int, url string) {
			// Decrement the counter when the goroutine completes.
			defer wg.Done()
			begin := time.Now()
			// Fetch the URL.
			resp, _ := http.Get(url)
			fmt.Printf("%d) Site %q took %s to retrieve with status code of %d.\n", i, url, time.Since(begin), resp.StatusCode)
		}(i, url)
	}
	// Wait for all HTTP fetches to complete.
	wg.Wait()
}
