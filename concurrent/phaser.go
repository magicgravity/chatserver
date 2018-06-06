package concurrent

import (
	"sync"
	"sync/atomic"
)

/*
 *  Phaser的高级用法，在Phaser内有2个重要状态，分别是phase和party。
 *  phase就是阶段，初值为0，当所有的线程执行完本轮任务，同时开始下一轮任务时，
 *  意味着当前阶段已结束，进入到下一阶段，phase的值自动加1。party就是线程，
 *  party=4就意味着Phaser对象当前管理着4个线程。Phaser还有一个重要的方法经常需要被重载，
 *  那就是boolean onAdvance(int phase, int registeredParties)方法。此方法有2个作用：
 *  1、当每一个阶段执行完毕，此方法会被自动调用，因此，重载此方法写入的代码会在每个阶段执行完毕时执行，
 *  相当于CyclicBarrier的barrierAction。
 *  2、当此方法返回true时，意味着Phaser被终止，因此可以巧妙的设置此方法的返回值来终止所有线程。
 */
type Phaser struct {
	phase int64
	party int64
	cond *sync.Cond
	onCallbackFunc *func()bool
	refCount int64
}

func  NewPhaser(parties int64,onAdvance *func()bool)*Phaser{
	ph := Phaser{}
	atomic.StoreInt64(&ph.party,parties)
	atomic.StoreInt64(&ph.phase,0)
	lock := sync.Mutex{}
	ph.cond = sync.NewCond(&lock)
	ph.onCallbackFunc= onAdvance
	atomic.StoreInt64(&ph.refCount,0)
	return &ph
}



func (ph *Phaser)ArriveAndAwaitAdvance(){
	atomic.AddInt64(&ph.refCount,1)
	if isAllArrive:=atomic.CompareAndSwapInt64(&ph.refCount,ph.party,ph.party);isAllArrive{
		if ph.onCallbackFunc!= nil {
			(*ph.onCallbackFunc)()
		}
		atomic.StoreInt64(&ph.refCount,0)
		ph.cond.L.Lock()
		defer ph.cond.L.Unlock()
		ph.cond.Broadcast()
	}else{
		ph.cond.L.Lock()
		defer ph.cond.L.Unlock()
		ph.cond.Wait()
	}
}