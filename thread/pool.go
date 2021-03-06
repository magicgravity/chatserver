package thread

import (
	"time"
	"sync"
	"math"
	"errors"
	"github.com/magicgravity/chatserver/common"
	"fmt"
	"log"
)

type RefusePolicy int

const (
	AbortPolicy RefusePolicy = 0		/* 丢弃任务并抛出RejectedExecutionError异常。 */
	DiscardPolicy RefusePolicy =1		/* 丢弃任务，但是不抛出异常 */
	DiscardOldestPolicy  RefusePolicy=2	/* 丢弃队列最前面的任务，然后重新尝试执行任务（重复此过程）*/
	CallerRunsPolicy RefusePolicy=3		/* 调用线程处理该任务*/

	PoolState_Inital = 0
	PoolState_Running = 1
	PoolState_Destroy = 2
)

var RejectedExecutionError  = errors.New("task is rejected to execute")
var RejectedOldestTaskExecutionError = errors.New("task is too old ,so it was rejected to execute")
var WarpFunctionToTaskError = errors.New("warp function to worktask fail")
var WorkTaskParamInValidCommonError = errors.New("worktask param is invalid ")



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
	submit(int,interface{},[]interface{})chan *TaskResult
	start()
}


type threadPoolExecutor struct {
	workState int			/* 0:初始状态  1:运行状态    */
	mainLock *sync.Mutex
	pool *goRoutinePool
	waitQueue []*WorkTask
	maxWaitQueueSize int
	refusePolicy RefusePolicy
	exitChan chan bool
	noticeChCh chan chan GNotice
}


func (pe *threadPoolExecutor)execute(job *WorkTask){
	var isGoHere = false

	log.Printf("---------execute----->%d \r\n",pe.workState)

	addToWaitQueue := func(j *WorkTask){
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
	if pe.workState == PoolState_Inital {
		//初始状态  先存到队列里
		addToWaitQueue(job)
		if isGoHere{
			goto here
		}
		return
	}else if pe.workState ==PoolState_Running{
		//运行状态
		if len(pe.waitQueue)>0 {
			log.Println("execute  ... addToWaitQueue")
			//仍然往等待队列里提交
			go addToWaitQueue(job)
			return
		}else{
			//直接提交到pool里
			log.Println("execute  ... addToWorkPool")
			pe.pool.workQueue<-job
			//go func(){
			//		pe.pool.workQueue<-job
			//	}()
		}
	}else{
		//todo 其他状态
	}
}

/*
	提交一个任务
	p：优先级
	function :任务描述
	params:	  任务参数

 */
func (pe *threadPoolExecutor)Submit(p int,function interface{},params []interface{})(chan *TaskResult,error){
	//TODO 需要判断类型 把提交的内容 包装成WorkTask
	//func([]interface{})([]interface{},error)
	task := WorkTask{}
	task.Priority = p
	task.Params = params
	futureChan := make(chan *TaskResult)
	task.FutureResult = futureChan
	var warpError error = nil
	switch function.(type) {
	case func(interface{})(interface{},error):
		if len(params)!=1{
			warpError = WorkTaskParamInValidCommonError
		}
		task.Tasks = func(p []interface{})([]interface{},error){
			func1 := function.(func(interface{})(interface{},error))
			rs,err := func1(p[0])
			rss := make([]interface{},1)
			rss[0] = rs
			return rss,err
		}
	case func([]interface{})([]interface{},error):
		task.Tasks = function.(func([]interface{})([]interface{},error))
	case func([]interface{}):
		task.Tasks = func(p []interface{}) ([]interface{}, error) {
			func1 := function.(func([]interface{}))
			func1(p)
			return nil,nil
		}
	case func(interface{}):
		if len(params)!=1{
			warpError = WorkTaskParamInValidCommonError
		}
		task.Tasks = func(p []interface{}) ([]interface{}, error) {
			func1 := function.(func(interface{}))
			func1(p)
			return nil,nil
		}
	case func(interface{})(error):
		if len(params)!=1{
			warpError = WorkTaskParamInValidCommonError
		}
		task.Tasks = func(p []interface{}) ([]interface{}, error) {
			func1 := function.(func(interface{})(error))
			return nil,func1(p)
		}
	case func([]interface{})(interface{},error):
		task.Tasks = func(p []interface{}) ([]interface{}, error) {
			func1 := function.(func([]interface{})(interface{},error))
			v,r := func1(p)
			rss := make([]interface{},1)
			rss[0] = v
			return rss,r
		}
	case func([]interface{})interface{}:
		task.Tasks = func(p []interface{}) ([]interface{}, error) {
			func1 := function.(func([]interface{})interface{})
			v := func1(p)
			rss := make([]interface{},1)
			rss[0] = v
			return rss,nil
		}
	default:
		warpError = WarpFunctionToTaskError

	}
	if warpError== nil {
		pe.execute(&task)
		return futureChan,warpError
	}else{
		return nil,warpError
	}
}

func (pe *threadPoolExecutor)shutDown(){
	pe.mainLock.Lock()
	defer pe.mainLock.Unlock()
	pe.workState = PoolState_Destroy
	pe.exitChan <- true
	endCmd := Op{GoRoutine_OpCmd_End}
	//通知池内的goroutine
	for _,g :=range pe.pool.pool{
		g.opChan <-endCmd
	}
}

/*
	启动
 */
func (pe *threadPoolExecutor)Start(){
	log.Printf("-------------Start  --> %d \r\n",pe.workState)
	if pe.workState == PoolState_Inital {
		//只有初始状态 才进入
		pe.mainLock.Lock()
		pe.workState = PoolState_Running
		defer pe.mainLock.Unlock()

		go func(){
			nclist := make([]chan GNotice,100)
			for{
				select {
				case nc:= <-pe.noticeChCh:
					nclist = append(nclist,nc)
				default:
					for _,n := range nclist{
						select {
						case notice:=<-n:
							fmt.Printf("goroutine notice ==> %v \r\n",notice)
						}
					}
				}
			}
		}()

		go func(){
		forExit:
			for {
				if len(pe.waitQueue) > 0 {
					top := pe.waitQueue[0]
					select {
					case pe.pool.workQueue <- top:
						//添加成功 则删除这个任务
						pe.mainLock.Lock()
						pe.waitQueue = pe.waitQueue[1:]
						pe.mainLock.Unlock()
					default:
						//TODO
					}
				}else if pe.workState == PoolState_Destroy{
					//销毁情况 则需要退出
					break forExit
				}else{
					time.Sleep(time.Second*5)
				}
			}
		}()

		go func(){
		forExit2:
			for {
				//log.Printf("~~~GoRoutine_Current poolSize => %d ,corePoolSize => %d \r\n",len(pe.pool.pool),pe.pool.corePoolSize)
				if len(pe.pool.pool)<pe.pool.corePoolSize{
					//线程数量小与核心池数量  则补充成员

					checkPoolSizeForCreate(pe)
				}

				//log.Printf("~~~GoRoutine_poolSize => %d \r\n",len(pe.pool.pool))
				//分配任务
				select {
				case job:= <- pe.pool.workQueue:
					log.Println("found one job !!!")
					var retryTime = 0
					ranMaxRetryTime := common.GenRandomInt(1,len(pe.pool.pool))

					for _,q := range pe.pool.pool{
						if q.state == GoRoutine_IntialStatus{
							log.Println("~~~GoRoutine_IntialStatus  >> dispatch")
							//如果是初始状态  启动它
							q.Run()
							q.Dispatch(job)
							//已分配的则 跳出for 处理
							break
						}else if q.state == GoRoutine_IdleStatus {
							log.Println("~~~GoRoutine_IdleStatus  >> dispatch")
							//
							q.Dispatch(job)
							//已分配的则 跳出for 处理
							break
						}else if q.state == GoRoutine_EndStatus {
							log.Println("~~~GoRoutine_EndStatus  >> dispatch")
							//如果是终止状态 检查数量 是否需要补足
							if len(pe.pool.pool)<pe.pool.corePoolSize{
								checkPoolSizeForCreate(pe)
							}
							continue
						}else if q.state == GoRoutine_WaitRunStatus  || q.state == GoRoutine_RunningStatus{
							log.Println("~~~GoRoutine_WaitRunStatus  >> dispatch")
							//如果是 等待运行或者 运行状态
							if retryTime >= ranMaxRetryTime{
								//超出次数 就往这个发
								q.Dispatch(job)
								break
							}else{
								//没有找到 就继续找
								continue
							}
						}
					}

					//TODO 有可能没找到合适的 分配任务

				case exitFlag:= <-pe.exitChan :
					if exitFlag {
						break forExit2
					}


				}
			}
		}()

	}
}



func checkPoolSizeForCreate(pe *threadPoolExecutor){
	for i:=0;i< pe.pool.corePoolSize;i++ {
		if g,ok := pe.pool.pool[i];ok {
			log.Println("checkPoolSizeForCreate ok=true")
			// 只有非终止状态的 不会启动新线程
			if g.state != GoRoutine_EndStatus {
				continue
			}
		}else{
			log.Println("checkPoolSizeForCreate ok=false")
			_,n :=createGroutineAndAddPool(pe.pool,GoRoutine_GType_Other,i)
			pe.noticeChCh<-n
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
func (p *executorsFactory)NewCachedThreadPool()(*threadPoolExecutor,error){
	tpool := threadPoolExecutor{}
	tpool.mainLock = &sync.Mutex{}
	tpool.workState = 0
	tpool.maxWaitQueueSize = 1000
	tpool.waitQueue = make([]*WorkTask,0)
	tpool.refusePolicy = DiscardOldestPolicy
	tpool.noticeChCh = make(chan chan GNotice)
	tpool.exitChan = make(chan bool)
	grpool := goRoutinePool{}
	grpool.corePoolSize = 0
	grpool.maxPoolSize = math.MaxInt32
	grpool.keepAliveTime = time.Minute
	grpool.pool = make(map[int]*GoRoutine)
	grpool.workQueue = make(chan *WorkTask,5)

	tpool.pool = &grpool

	return &tpool,nil
}

/*
创建一个定长线程池，可控制线程最大并发数，超出的线程会在队列中等待。
 */
func (p *executorsFactory)NewFixedThreadPool(nThread int)(*threadPoolExecutor,error){
	tpool := threadPoolExecutor{}
	tpool.mainLock = &sync.Mutex{}
	tpool.workState = 0
	tpool.maxWaitQueueSize = nThread
	tpool.waitQueue = make([]*WorkTask,0)
	tpool.refusePolicy = AbortPolicy
	tpool.noticeChCh = make(chan chan GNotice)
	tpool.exitChan = make(chan bool)
	grpool := goRoutinePool{}
	grpool.corePoolSize = nThread
	grpool.maxPoolSize = math.MaxInt32
	grpool.keepAliveTime = time.Minute
	grpool.pool = make(map[int]*GoRoutine,nThread)
	grpool.workQueue = make(chan *WorkTask,5)

	for i:=0;i<nThread;i++{
		if i==0 {
			createGroutineAndAddPool(&grpool,GoRoutine_GType_Main,i)
		}else{
			createGroutineAndAddPool(&grpool,GoRoutine_GType_Other,i)
		}

	}
	tpool.pool = &grpool

	return &tpool,nil
}


func createGroutineAndAddPool(p *goRoutinePool,gtype,pos int)(chan Op,chan GNotice){
	//log.Printf("=============>createGroutineAndAddPool===pool len >>>%d , pos ===>>> %d \r\n",len(p.pool),pos)

	opChannel := make(chan Op)
	noticeChannel := make(chan GNotice)
	p.pool[pos] = makeNewGoRoutine(pos,gtype,GoRoutine_DefaultQueueSize,
		GoRoutine_DefaultAliveTime,opChannel,noticeChannel)
	return opChannel,noticeChannel

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