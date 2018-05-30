package concurrent

import (
	"sync"
	"sync/atomic"
	"fmt"
)

type CyclicBarrier struct {
	num int64
	refCount int64
	cond *sync.Cond
	barrierFuc *func()bool
	isBroken int32
}


func NewCyclicBarrier(parties int,barrierAction *func()bool)*CyclicBarrier{
	cb := CyclicBarrier{}
	cb.num = int64(parties)
	atomic.StoreInt64(&cb.refCount,0)
	lock := sync.Mutex{}
	cb.cond = sync.NewCond(&lock)
	cb.barrierFuc =  barrierAction
	atomic.StoreInt32(&cb.isBroken,0)
	return &cb
}


func (cb *CyclicBarrier)Await(){
	if isBroken :=atomic.CompareAndSwapInt32(&cb.isBroken,1,1);isBroken {
		fmt.Println("broken state")
		return
	}

	atomic.AddInt64(&cb.refCount,1)
	if isSwap := atomic.CompareAndSwapInt64(&cb.refCount,cb.num,cb.num);isSwap{
		fmt.Println("reach the limit")
		cb.cond.L.Lock()
		defer cb.cond.L.Unlock()
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("barrierFuc exec fail >>：", err)
				atomic.StoreInt32(&cb.isBroken,1)
			}
		}()
		var isOk = true
		if cb.barrierFuc != nil{
			isOk = (*cb.barrierFuc)()
		}
		if !isOk{
			atomic.StoreInt32(&cb.isBroken,1)
		}
		atomic.StoreInt64(&cb.refCount, 0)

		cb.cond.Broadcast()

	}else{
		cb.cond.L.Lock()
		defer cb.cond.L.Unlock()
		cb.cond.Wait()
	}
}


func (cb *CyclicBarrier)Reset(){
	if isBroken :=atomic.CompareAndSwapInt32(&cb.isBroken,1,1);isBroken {
		//broken state
		return
	}

	cb.cond.L.Lock()
	defer cb.cond.L.Unlock()
	atomic.StoreInt64(&cb.refCount,0)
	cb.cond.Broadcast()
}


