package main

import (
	"sync"
)

// 使用slice和条件变量实现一个简单的多生产者多消费者队列
// https://github.com/Chasiny/Blog/blob/master/blog/go/%E4%BD%BF%E7%94%A8slice%E5%92%8C%E6%9D%A1%E4%BB%B6%E5%8F%98%E9%87%8F%E5%AE%9E%E7%8E%B0%E4%B8%80%E4%B8%AA%E7%AE%80%E5%8D%95%E7%9A%84%E5%A4%9A%E7%94%9F%E4%BA%A7%E8%80%85%E5%A4%9A%E6%B6%88%E8%B4%B9%E8%80%85%E9%98%9F%E5%88%97.md

type MessageQueue struct {
	msgdata []interface{} // 缓冲区pipeline
	len     int32         // 缓冲区长度

	readPos   int32      // 读到指向的指针
	readMutex sync.Mutex // 读取锁

	writePos   int32      // 写入指向的指针
	wirteMutex sync.Mutex // 写入锁

	emptyCond *sync.Cond // 缓冲区为空条件变量
	fullCond  *sync.Cond // 缓冲区为满条件变量
}

// 写入的方法(Put)
func (mq *MessageQueue) Put(in interface{}) {
	// 首先获取写锁, 所有写入的优先级是一样的
	mq.wirteMutex.Lock()
	defer mq.wirteMutex.Unlock()

	// 判断缓冲区是否为满
	mq.fullCond.L.Lock()
	defer mq.fullCond.L.Unlock()
	for (mq.writePos+1)%mq.len == mq.readPos {
		// 缓冲区为满, 等待消费者消费的通知缓冲区有数据被取出
		mq.fullCond.Wait()
	}

	// 写入一个数据
	mq.msgdata[mq.writePos] = in
	mq.writePos = (mq.writePos + 1) % mq.len

	// 通知消费者缓冲区已经有数据了
	mq.emptyCond.Signal()
}

// 读取的方法(Get)
func (mq *MessageQueue) Get(out interface{}) {
	// 获取读锁, 读取的优先级也是一样的
	mq.readMutex.Lock()
	defer mq.readMutex.Unlock()

	// 判断缓冲区是否为空
	mq.emptyCond.L.Lock()
	defer mq.emptyCond.L.Unlock()

	for mq.writePos == mq.readPos {
		// 缓冲区为空, 等待生产者通知缓冲区有数据存入
		mq.emptyCond.Wait()
	}

	// 读取
	out = mq.msgdata[(mq.readPos)%mq.len]
	mq.readPos = (mq.readPos + 1) % mq.len

	// 通知生产者已经有缓冲区有空间了
	mq.fullCond.Signal()

	return
}

// 长度(Len)
func (mq *MessageQueue) Len() int32 {
	if mq.writePos < mq.readPos {
		return mq.writePos + mq.len - mq.readPos
	}

	return mq.writePos - mq.readPos
}

// New方法
func NewMQ(len int32) *MessageQueue {
	if len < 1 {
		panic("new meg queue fail: len < 1")
		return nil
	}
	l := &sync.Mutex{}
	return &MessageQueue{
		msgdata:  make([]interface{}, len+1),
		len:      len + 1,
		readPos:  0,
		writePos: 0,

		emptyCond: sync.NewCond(l),
		fullCond:  sync.NewCond(l),
	}
}

func main() {

}
