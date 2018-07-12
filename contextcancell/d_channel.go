package main

import (
	"log"
	"sync/atomic"
)

// 参考 https://blog.csdn.net/qq_26981997/article/details/73252886

type (
	Message1 struct{}
	Message2 struct{}
)

type A struct {
	close  int32
	msgbuf chan interface{}
}

func NewA() *A {
	a := &A{
		msgbuf: make(chan interface{}, 10),
	}
	go a.receive()
	return a
}

func (a *A) Post(message interface{}) {
	if atomic.LoadInt32(&a.close) == 1 {
		a.msgbuf <- message
	}
}

func (a *A) Close() {
	if atomic.CompareAndSwapInt32(&a.close, 0, 1) {
		// do other thing
		close(a.msgbuf)
	}
}

func (a *A) receive() {
	// 通过defer 实现简单的故障隔离
	defer func() {
		if er := recover(); er != nil {
			log.Println(er)
		}
	}()
	// 执行消息处理
	for message := range a.msgbuf {
		switch msg := message.(type) {
		case Message1:
			a.foo1(msg)
		case Message2:
			a.foo2(msg)
		}
	}
}

func (a *A) foo1(msg Message1) {

}

func (a *A) foo2(msg Message2) {

}
