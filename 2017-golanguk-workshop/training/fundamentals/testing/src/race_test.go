package main

import (
	"fmt"
	"sync"
	"testing"
)

var m = 0

func inc() {
	m++
}

func TestRace(t *testing.T) {
	wg := &sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			inc()
		}()
		go func() {
			defer wg.Done()
			fmt.Println(m)
		}()
	}
	wg.Wait()
	if m == 0 {
		t.Fatalf("expect m (%d) to not be zero", m)
	}
}
