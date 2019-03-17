package main

import (
	"fmt"
	"github.com/petar/GoLLRB/llrb"
	"math/rand"
	"time"
)

// 优先级队列
// 参考: https://gist.github.com/fmstephe/4fdc930ff180be3e92c693ad5a24e1b3

func main() {
	q := NewPrioq()
	s := time.Now()
	go put(200000, q)
	go put(200000, q)
	go put(200000, q)
	get(200000, q)
	get(200000, q)
	get(200000, q)

	fmt.Println(time.Now().Sub(s))
}

type timestamp int64

func (t timestamp) Less(i llrb.Item) bool {
	return t < i.(timestamp)
}

func put(count int, q *Prioq) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < count; i++ {
		q.In <- timestamp(r.Int63())
	}
}

func get(count int, q *Prioq) {
	for i := 0; i < count; i++ {
		<-q.Out
	}
}

type Prioq struct {
	In   chan llrb.Item
	tree *llrb.LLRB
	Out  chan llrb.Item
}

func NewPrioq() *Prioq {
	q := &Prioq{
		In:   make(chan llrb.Item, 100),
		tree: llrb.New(),
		Out:  make(chan llrb.Item, 100),
	}
	q.run()
	return q
}

func (q *Prioq) run() {
	go func() {
		for {
			i := q.getIn()
			if i != nil {
				q.tree.ReplaceOrInsert(i)
			}
			o := q.tree.Max()
			if o != nil {
				q.putOut(o)
			}
		}
	}()
}

func (q *Prioq) getIn() llrb.Item {
	select {
	case i := <-q.In:
		return i
	default:
		return nil
	}
}

func (q *Prioq) putOut(i llrb.Item) {
	q.Out <- i
}
