package thread

import (
	"time"
	"sync"
	"math"
	"errors"
)

type RefusePolicy int

const (
	AbortPolicy RefusePolicy = 0		/* 丢弃任务并抛出RejectedExecutionError异常。 */
	DiscardPolicy RefusePolicy =1		/* 丢弃任务，但是不抛出异常 */
	DiscardOldestPolicy  RefusePolicy=2	/* 丢弃队列最前面的任务，然后重新尝试执行任务（重复此过程）*/
	CallerRunsPolicy RefusePolicy=3		/* 调用线程处理该任务*/
)

var RejectedExecutionError  = errors.New("task is rejected to execute")
var RejectedOldestTaskExecutionError = errors.New("task is too old ,so it was rejected to execute")



type goRoutinePool struct {
	maxPoolSize int
	corePoolSize int
	keepAliveTime time.Duration
	pool map[int]*GoRoutine
	workQueue chan *WorkTask
}

type Executor interface {
	execute(*WorkTask)
}

type ExecutorService interface {
	shutDown()
	submit(*WorkTask)
	start()
}


type threadPoolExecutor struct {
	workState int			/* 0:初始状态  1:运行状态    */
	mainLock *sync.Mutex
	pool *goRoutinePool
	waitQueue []*WorkTask
	maxWaitQueueSize int
	refusePolicy RefusePolicy
}


func (pe *threadPoolExecutor)execute(job *WorkTask){

	var isGoHere = false

	addToWaitQueue := func(){
		isGoHere = false
		if len(pe.waitQueue)<pe.maxWaitQueueSize{
			pe.mainLock.Lock()
			pe.waitQueue = append(pe.waitQueue,job)
			pe.mainLock.Unlock()
			return
		}else{
			//队列满了
			switch pe.refusePolicy {
			case AbortPolicy:
				rst := TaskResult{}
				rst.err = RejectedExecutionError
				rst.data = nil
				job.FutureResult<-&rst
				close(job.FutureResult)
			case DiscardPolicy:
				rst := TaskResult{}
				rst.err = nil
				rst.data = nil
				job.FutureResult<-&rst
				close(job.FutureResult)
				return
			case DiscardOldestPolicy:
				oldest := pe.waitQueue[0]
				pe.mainLock.Lock()
				pe.waitQueue = pe.waitQueue[1:]
				pe.mainLock.Unlock()
				rst := TaskResult{}
				rst.err = RejectedOldestTaskExecutionError
				rst.data = nil
				oldest.FutureResult<-&rst
				close(oldest.FutureResult)
				isGoHere = true
			case CallerRunsPolicy:
			default:
				//默认直接丢弃  不报错
				rst := TaskResult{}
				rst.err = nil
				rst.data = nil
				job.FutureResult<-&rst
				close(job.FutureResult)
				return
			}
		}
	}


here:
	if pe.workState == 0 {
		//初始状态  先存到队列里
		addToWaitQueue()
		if isGoHere{
			goto here
		}
		return
	}else if pe.workState ==1{
		//运行状态
		if len(pe.waitQueue)>0 {
			//仍然往等待队列里提交
			addToWaitQueue()
			return
		}else{
			//直接提交到pool里
			go func(){
					pe.pool.workQueue<-job
				}()
		}
	}else{
		//todo 其他状态
	}
}

func (pe *threadPoolExecutor)submit(job *WorkTask){

}

/*
	启动
 */
func (pe *threadPoolExecutor)start(){
	if pe.workState == 0 {
		//只有初始状态 才进入
		pe.mainLock.Lock()
		defer pe.mainLock.Unlock()
		size := len(pe.pool.pool)
		var idx =0
		for {
			if idx <size{


			}else{
				break
			}
		}
	}
}

type executorsFactory struct {

}

var executors *executorsFactory
var once sync.Once

func GetExecutors()*executorsFactory{
	once.Do(func() {
		executors = &executorsFactory{}
	})
	return executors
}


/*
创建一个可缓存线程池，如果线程池长度超过处理需要，可灵活回收空闲线程，若无可回收，则新建线程
 */
func (p *executorsFactory)newCachedThreadPool()(*threadPoolExecutor,error){
	tpool := threadPoolExecutor{}
	tpool.mainLock = &sync.Mutex{}
	tpool.workState = 0
	tpool.maxWaitQueueSize = 1000
	tpool.waitQueue = make([]*WorkTask,0)
	tpool.refusePolicy = DiscardOldestPolicy
	grpool := goRoutinePool{}
	grpool.corePoolSize = 0
	grpool.maxPoolSize = math.MaxInt32
	grpool.keepAliveTime = time.Minute
	grpool.pool = make(map[int]*GoRoutine)

	tpool.pool = &grpool

	return &tpool,nil
}

/*
创建一个定长线程池，可控制线程最大并发数，超出的线程会在队列中等待。
 */
func (p *executorsFactory)newFixedThreadPool(nThread int)(*threadPoolExecutor,error){
	tpool := threadPoolExecutor{}
	tpool.mainLock = &sync.Mutex{}
	tpool.workState = 0
	tpool.maxWaitQueueSize = nThread
	tpool.waitQueue = make([]*WorkTask,0)
	tpool.refusePolicy = AbortPolicy
	grpool := goRoutinePool{}
	grpool.corePoolSize = nThread
	grpool.maxPoolSize = math.MaxInt32
	grpool.keepAliveTime = time.Minute
	grpool.pool = make(map[int]*GoRoutine,nThread)

	tpool.pool = &grpool

	return &tpool,nil
}

/*
创建一个定长线程池，支持定时及周期性任务执行。
 */
func (p *executorsFactory)newScheduledThreadPool(){

}


/*
创建一个单线程化的线程池，它只会用唯一的工作线程来执行任务，保证所有任务按照指定顺序(FIFO, LIFO, 优先级)执行。
 */
func (p *executorsFactory)newSingleThreadExecutor(){

}