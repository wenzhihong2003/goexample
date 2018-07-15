package gp

import (
	"sync"
	"sync/atomic"
	"time"
)

// Pool is a struct to represent goroutine pool.
type Pool struct {
	stack       []*goroutine
	idleTimeout time.Duration
	sync.Mutex
}

// goroutine is actually a background goroutine, with a channel binded for communication.
type goroutine struct {
	ch     chan func()
	status int32
}

const (
	statusIdle  int32 = 0
	statusInUse int32 = 1
	statusDead  int32 = 2
)

// New returns a new *Pool object
func New(idleTimeout time.Duration) *Pool {
	pool := &Pool{
		idleTimeout: idleTimeout,
		stack:       make([]*goroutine, 0, 64),
	}
	return pool
}

// Go works like go func(), but goroutines are pooled for reusing.
// This strategy can avoid runtime.morestack, because pooled goroutine is already enlarged.
func (pool *Pool) Go(f func()) {
	for {
		g := pool.get()
		if atomic.CompareAndSwapInt32(&g.status, statusIdle, statusInUse) {
			g.ch <- f
			return
		}
		// Status already changed from statusIdle => statusDead, drop it, find next one.
	}
}

func (pool *Pool) get() *goroutine {
	pool.Lock()
	if len(pool.stack) == 0 {
		pool.Unlock()
		return pool.alloc()
	}
	ret := pool.stack[len(pool.stack)-1]
	pool.stack = pool.stack[:len(pool.stack)-1]
	pool.Unlock()
	return ret
}

func (pool *Pool) alloc() *goroutine {
	g := &goroutine{
		ch: make(chan func()),
	}
	go g.workLoop(pool)
	return g
}

func (g *goroutine) put(pool *Pool) {
	g.status = statusInUse
	pool.Unlock()

	// Recycle dead goroutine space.
	i := 0
	for ; i < len(pool.stack) && atomic.LoadInt32(&pool.stack[i].status) == statusDead; i++ {
	}
	pool.stack = append(pool.stack[i:], g)
	pool.Unlock()
}

func (g *goroutine) workLoop(pool *Pool) {
	timer := time.NewTimer(pool.idleTimeout)
	for {
		select {
		case <-timer.C:
			// Check to avoid a corner case that the goroutine is take out from pool,
			// and get this signal at the same time.
			succ := atomic.CompareAndSwapInt32(&g.status, statusIdle, statusDead)
			if succ {
				return
			}
		case work := <-g.ch:
			work()
			// Put g back to the pool.
			// This is the normal usage for a resource pool:
			//
			//     obj := pool.get()
			//     use(obj)
			//     pool.put(obj)
			//
			// But when goroutine is used as a resource, we can't pool.put() immediately,
			// because the resource(goroutine) maybe still in use.
			// So, put back resource is done here,  when the goroutine finish its work.
			g.put(pool)
		}
		timer.Reset(pool.idleTimeout)
	}
}