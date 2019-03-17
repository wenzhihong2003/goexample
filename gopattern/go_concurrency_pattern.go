package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// go的并发模式. 参考:https://medium.com/@thejasbabu/concurrency-patterns-golang-5c5e1bcd0833

// 生成器模式

func fib1(n int) <-chan int {
	c := make(chan int)
	go func() {
		for i, j := 0, 1; i < n; i, j = i+j, i {
			c <- i
		}
		close(c)
	}()
	return c
}

func fib1Main() {
	// fib returns the fibonacci numbers lesser than 1000
	for i := range fib1(1000) {
		// consumer which consumes the data produced by the generator, which further does some extra computations
		v := i * i
		fmt.Println(v)
	}
}

// future 模式, 类似于jquery的future/promise模式
type data struct {
	Body  []byte
	Error error
}

func futureData(url string) <-chan data {
	c := make(chan data, 1)

	go func() {
		var body []byte
		var err error
		resp, err := http.Get(url)
		defer resp.Body.Close()

		body, err = ioutil.ReadAll(resp.Body)
		c <- data{Body: body, Error: err}
	}()

	return c
}

func futureDataMain() {
	future := futureData("http://test.future.com")

	// do many other things

	body := <-future
	fmt.Printf("response:%#v", string(body.Body))
	fmt.Printf("error:%#v", body.Error)
}

// fan-in fan-out

func generatePipeline(numbers []int) <-chan int {
	out := make(chan int)

	go func() {
		for _, n := range numbers {
			out <- n
		}
		close(out)
	}()

	return out
}

func squareNumber(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()

	return out
}

func fanIn(input1, input2 <-chan int) <-chan int {
	c := make(chan int)

	go func() {
		for {
			select {
			case s := <-input1:
				c <- s
			case s := <-input2:
				c <- s
			}
		}
	}()

	return c
}

func faninoutMain() {
	randomNumbers := []int{13, 44, 56, 99, 9, 45, 67, 90, 78, 23}
	// generate the common channel with inputs
	inputChan := generatePipeline(randomNumbers)

	// fan-out to 2 Go-routine
	c1 := squareNumber(inputChan)
	c2 := squareNumber(inputChan)

	// Fan-in the resulting squared numbers
	c := fanIn(c1, c2)
	sum := 0

	// do the summation
	for i := 0; i < len(randomNumbers); i++ {
		sum += <-c
	}
	fmt.Printf("Total Sum of Squares: %d", sum)
}
