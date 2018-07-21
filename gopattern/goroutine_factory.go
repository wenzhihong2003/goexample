package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

// 通常使用工厂方法将goroutine和通道绑定。
// https://www.jianshu.com/p/1132baa9475c

type receiver struct {
	wg   sync.WaitGroup
	data chan int
}

func newReceiver() *receiver {
	r := &receiver{
		data: make(chan int),
	}
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		for x := range r.data {
			fmt.Println("recv:", x)
		}
	}()
	return r
}

func main() {
	r := newReceiver()
	r.data <- 1
	r.data <- 2
	close(r.data) // 关闭通道，发出结束通知
	r.wg.Wait()   // 等待接收者处理结束
}

// 鉴于通道本身就是一个并发安全的队列，可用作ID generator、Pool等用途
type pool chan []byte

func newPool(cap int) pool {
	return make(chan []byte, cap)
}

func (p pool) get() []byte {
	var v []byte
	select {
	case v = <-p:
	default:
		v = make([]byte, 10) // 返回失败，新建
	}
	return v
}

func (p pool) put(b []byte) {
	select {
	case p <- b: // 放回
	default: // 放回失败, 新建

	}
}

// 用通道实现信号量（semaphore）
func semaphore() {
	runtime.GOMAXPROCS(4)
	var wg sync.WaitGroup

	sem := make(chan struct{}, 2) // 最多允许两个并发同时执行
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sem <- struct{}{}        // acquire: 获取信号
			defer func() { <-sem }() // release: 释放信号
			time.Sleep(time.Second * 2)
			fmt.Println(id, time.Now())
		}(i)
	}

	wg.Wait()
}

// atexit函数是一个特殊的函数，它是在正常程序退出时调用的函数，我们把他叫为登记函数（函数原型：int atexit (void (*)(void))）：
// ⼀个进程可以登记若⼲个（具体⾃⼰验证⼀下）个函数，这些函数由exit⾃动调⽤，这些函数被称为终⽌处理函数， atexit函数可以登记这些函数。 exit调⽤终⽌处理函数的顺序和atexit登记的顺序相反（网上很多说造成顺序相反的原因是参数压栈造成的，参数的压栈是先进后出，和函数的栈帧相同），如果⼀个函数被多次登记，也会被多次调⽤。
// python中有专门的atexit模块，简介如下：
// 从模块的名字也可以看出来，atexit模块主要的作用就是在程序即将结束之前执行的代码，atexit模块使用register函数用于注册程序退出时的回调函数，然后在回调函数中做一些资源清理的操作。
// 注意：
// 1，如果程序是非正常crash，或通过os._exit()退出，注册的回调函数将不会被调用。
// 2，也可以通过sys.exitfunc来注册回调，但通过它只能注册一个回调，而且还不支持参数。
// 3，建议使用atexit来注册回调函数。
//

var exits = &struct {
	sync.RWMutex
	signals chan os.Signal
	funcs   []func()
}{}

func atexit(f func()) {
	exits.Lock()
	defer exits.Unlock()
	exits.funcs = append(exits.funcs, f)
}

func waitExit() {
	if exits.signals == nil {
		exits.signals = make(chan os.Signal)
		signal.Notify(exits.signals, syscall.SIGINT, syscall.SIGTERM)
		fmt.Println("test")
	}
	exits.RLock()
	for _, f := range exits.funcs {
		defer f() // 延迟调用函数采用FILO顺序执行。即便某些函数panic，延迟调用也能确保后续函数执行。
	}
	fmt.Println("after range exits.funcs")
	exits.RUnlock()
	fmt.Println("after exits.Runlock")
	<-exits.signals
}

func examAtexit() {
	atexit(func() {
		fmt.Println("exit1....")
	})
	atexit(func() {
		fmt.Println("exit2...")
	})
	fmt.Println("before exit")
	waitExit()

}
