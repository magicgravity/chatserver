package concurrent

import (
	"sync/atomic"
	"sync"
)

type CountDownLatch struct {
	num int
	chs chan bool
	refCount int64
	cond *sync.Cond
}

func NewCountDownLatch(count int)*CountDownLatch{
	cdl := CountDownLatch{}
	cdl.num = count
	cdl.chs =  make(chan bool,count)
	k := sync.Mutex{}
	cdl.cond = sync.NewCond(&k)
	atomic.StoreInt64(&cdl.refCount,int64(cdl.num))
	return &cdl
}

func (c *CountDownLatch)CountDown(){
	atomic.AddInt64(&c.refCount,-1)
	if isSwap := atomic.CompareAndSwapInt64(&c.refCount,0,0);isSwap {
		c.cond.L.Lock()
		defer c.cond.L.Unlock()
		c.cond.Broadcast()
	}
}

func (c *CountDownLatch)Await(){
	c.cond.L.Lock()
	defer c.cond.L.Unlock()
	c.cond.Wait()
}