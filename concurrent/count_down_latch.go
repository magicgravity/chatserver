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
	atomic.StoreInt64(&cdl.refCount,0)
	return &cdl
}

func (c *CountDownLatch)CountDown(){
	//c.chs<-true
	atomic.AddInt64(&c.refCount,-1)
	c.cond.L.Lock()
	defer c.cond.L.Unlock()
	rc := atomic.LoadInt64(&c.refCount)
	if rc==0{
		c.cond.Broadcast()
	}
}

func (c *CountDownLatch)Await(){
	isSwap := atomic.CompareAndSwapInt64(&c.refCount,int64(c.num),int64(c.num))
	if isSwap{
		return
	}
	atomic.AddInt64(&c.refCount,1)
	c.cond.L.Lock()
	defer c.cond.L.Unlock()
	c.cond.Wait()
}