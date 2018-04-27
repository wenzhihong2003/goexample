package main

import (
	"time"
	"math/rand"
	"fmt"
	"log"
)

// 参考: https://medium.com/justforfunc/why-are-there-nil-channels-in-go-9877cc0b2308

func asChain(vs ...int) <-chan int {
	c := make(chan int)
	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Microsecond)
		}
		close(c)
	}()
	return c
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		defer close(c)
		for a != nil || b != nil {
			select {
			case v, ok := <-a:
				if !ok {
					log.Println("a is done")
					a = nil
					continue
				}
				c <- v
			case v, ok := <-b:
				if !ok {
					log.Println("b is done")
					b = nil
					continue
				}
				c <- v
			}
		}
	}()

	return c
}

func main() {
	a := asChain(1, 3, 5, 7)
	b := asChain(2, 4, 6, 8)
	c := merge(a, b)
	for v := range c {
		fmt.Println(v)
	}
}