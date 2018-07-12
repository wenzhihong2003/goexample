package main

import (
	"fmt"
	"time"
)

func main() {
	// 执行任务单元
	doWork := func(done <-chan interface{}, strings <-chan string) <-chan interface{} {
		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(terminated)
			for {
				select {
				case s := <-strings:
					// do something  这里执行相应的操作
					fmt.Println(s)
				case <-done:
					return
				}
			}

		}()

		return terminated
	}

	done := make(chan interface{})
	terminated := doWork(done, nil)

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling dowork goroutine...")
		close(done)
	}()

	<-terminated
	fmt.Println("done.")

}
