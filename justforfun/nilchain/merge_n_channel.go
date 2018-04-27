package main

// 参考: https://medium.com/justforfunc/two-ways-of-merging-n-channels-in-go-43c0b57cd1de

import (
	"time"
	"math/rand"
	"sync"
	"fmt"
	"reflect"
)

func asChain2(vs ...int) <-chan int {
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

func merge2(cs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	wg.Add(len(cs))
	for _, c := range cs {
		go func(c <-chan int) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func method1()  {
	a := asChain2(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
	b := asChain2(10, 11, 12, 13, 14, 15, 16, 17, 18, 19)
	c := asChain2(20, 21, 22, 23, 24, 25, 26, 27, 28, 29)
	for v := range merge2(a, b, c) {
		fmt.Println(v)
	}
}

func method2() {
	a := asChain2(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
	b := asChain2(10, 11, 12, 13, 14, 15, 16, 17, 18, 19)
	c := asChain2(20, 21, 22, 23, 24, 25, 26, 27, 28, 29)
	chans := []<-chan int {a, b, c}
	var cases []reflect.SelectCase

	// todo 这里用的第二种方法不成功, 先看看 reflect.SelectCase, reflect.Select怎么样处理
	for _, c := range chans {
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(c),
		})



	}



}

func main() {
	// method1()


}
