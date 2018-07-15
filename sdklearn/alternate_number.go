package main

import (
	"fmt"
	"io"
	"strconv"
)

// Golang让协程交替输出
// https://my.oschina.net/90design/blog/1609453

// 两个协程交替输出1-20
// 首先给通道A一个缓存，并在主进程中发送数据，使其堵塞，在第一个Goroutine中通道A接收并开始执行， 此时B是堵塞等待的， 等A执行完成发送数据到通道B， B开始执行。
func m1() {
	A := make(chan bool, 1)
	B := make(chan bool)
	Exit := make(chan bool)

	go func() {
		for i := 1; i <= 20; i++ {
			if ok := <-A; ok {
				fmt.Println("A=", 2*i-1)
				B <- true
			}
		}
	}()

	go func() {
		defer func() {
			close(Exit)
		}()
		for i := 1; i < 20; i++ {
			if ok := <-B; ok {
				fmt.Println("B: ", 2*i)
				A <- true
			}
		}
	}()
	A <- true
	<-Exit
}

// 使用一个chan搞定
func m2() {
	ch := make(chan int)
	exit := make(chan struct{})

	go func() {
		for i := 1; i <= 20; i++ {
			fmt.Println("g1:", <-ch)
			i++
			ch <- i
		}
	}()

	go func() {
		defer func() {
			close(ch)
			close(exit)
		}()
		for i := 0; i < 20; i++ {
			i++
			ch <- i
			fmt.Println("g2:", <-ch)
		}
	}()
	<-exit
}

var (
	num   int // 要输出的最大值
	line  = 0 // 通道发送计数器
	exit  = make(chan bool)
	chans []chan int // 要初始化的协程数量
)

func ChanWork(c chan int) {
	// 协程数
	lens := len(chans)
	for {
		// count 为输出计数器
		if count := <-chans[line]; count <= num {
			fmt.Println("channel "+strconv.Itoa(line)+" -> ", count)
			count++

			// 下一个发送通道
			line++
			if line >= lens {
				line = 0 // 循环，防止索引越界
			}
			go ChanWork(chans[line])
			chans[line] <- count
		} else {
			// 通道编号问题处理
			id := 0
			if line == 0 {
				id = lens - 1
			} else {
				id = line - 1
			}
			fmt.Println("在通道" + strconv.Itoa(id) + "执行完成")
			close(exit)
			return
		}
	}
}

// 在多个通道上完成
func m3() {

	chans = []chan int{
		make(chan int),
		make(chan int),
		make(chan int),
		make(chan int),
	}

	// 多协程启动入口
	go ChanWork(chans[0])
	// 接收要输出的最大数
	fmt.Println("输入要输出的最大数值:")
	_, ok := fmt.Scanf("%d\n", &num)
	if ok == io.EOF {
		return
	}
	// 触发协程同步执行
	chans[0] <- 1

	// 执行结束
	if <-exit {
		return
	}

}

func main() {
	m2()
}
