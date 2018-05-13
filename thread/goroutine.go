package thread

import (
	"sync"
	"time"
	"strconv"
	"sort"
)

const (
	GoRoutine_IntialStatus = 0
	GoRoutine_WaitRunStatus = 1
	GoRoutine_RunningStatus = 2
	GoRoutine_IdleStatus = 3
	GoRoutine_EndStatus = 4

	GoRoutine_OpCmd_End = "End"
	GoRoutine_OpCmd_Pause = "Pause"
	GoRoutine_OpCmd_Resume = "Resume"
)

type GoRoutineContext struct {
	dataMap map[string]interface{}
}

type TaskResult struct {
	data []interface{}
	err error
}

func MakeTaskResult(d []interface{},er error)*TaskResult{
	return &TaskResult{d,er}
}

type WorkTask struct {
	Tasks func([]interface{})([]interface{},error)
	Params []interface{}
	FutureResult chan *TaskResult
	Priority int
}

type TaskSet []*WorkTask

func (ts TaskSet)Len()int{
	return len(ts)
}

func (ts TaskSet)Swap(i,j int){
	ts[i],ts[j] = ts[j],ts[i]
}

/*
按从大到小排序
 */
func (ts TaskSet)Less(i, j int) bool{
	return ts[i].Priority>ts[j].Priority
}


type Handler func([]interface{})(bool,[]interface{})

type InnerTaskQueue struct {
	onBefore []Handler
	onAfter []Handler
	queue chan TaskSet
}

type Op struct {
	opCode string		/* 操作码 */
}

type GNotice struct {
	CurrentRunTask string
	LastRunResult []interface{}
	RunErr error
}

type GoRoutine struct {
	id int						/* 内部编号 */
	gtype int					/*	类型  0:主   1：其他普通	*/
	state int					/*  状态 0:初始状态  1：等待运行状态   2:运行状态   3：空闲状态   4:终止状态 */
	groupName string			/*	*/
	context *GoRoutineContext		/*  执行上下文  */
	localData map[string]interface{}		/* goroutine 内的本地空间 */
	lock *sync.Mutex
	rlock *sync.RWMutex
	Name string					/*名称 */
	Tasks *InnerTaskQueue		/*任务列表 */
	ctime time.Time				/*创建时间 */
	rtime time.Duration			/*运行时间	*/
	maxAliveTime time.Duration	/*指定的最大生存时间限制 超出自动销毁*/
	opChan chan Op				/*用于控制运行*/
	noticeChan chan GNotice		/*用于通知外部 */
	queueSize int				/* InnerTaskQueue 的queue 大小*/
	historyJobs TaskSet		/*已执行完成的历史任务 */
	tempTaskZone TaskSet  /*暂存的临时区域*/
}

var G_Group_Map map[string]*GGroup		/* 线程组 关系表 */

type GGroup struct {
	Name string							/* 名称 */
	Groutines *[]GoRoutine				/* 所有的线程 */
}

/*
	创建一个GoRoutine
 */
func makeNewGoRoutine(id,ctype,qsize int,maxATime time.Duration,opCh chan Op,noCh chan GNotice)*GoRoutine  {
	g :=GoRoutine{}
	g.id =  id
	g.queueSize = qsize
	g.Name = "groutine_"+strconv.Itoa(id)
	g.ctime = time.Now()
	g.state = GoRoutine_IntialStatus
	g.gtype = ctype
	g.maxAliveTime = maxATime
	g.opChan = opCh
	g.noticeChan = noCh
	g.localData =  make(map[string]interface{})
	g.lock = &sync.Mutex{}
	g.rlock = &sync.RWMutex{}
	inQ := InnerTaskQueue{}
	inQ.queue =  make(chan TaskSet,g.queueSize)
	inQ.onAfter = make([]Handler,0)
	inQ.onBefore = make([]Handler,0)
	g.Tasks = &inQ
	return &g
}
/*
	绑定一个上下文
 */
func (g *GoRoutine)Bind(ctx *GoRoutineContext)*GoRoutine{
	g.lock.Lock()
	if g.state ==GoRoutine_IntialStatus {
		g.context = ctx
	}
	g.lock.Unlock()
	return g
}


/*

 */
func (g *GoRoutine)AddHandlerToChain(htype bool,h Handler)*GoRoutine{
	g.lock.Lock()
	defer g.lock.Unlock()
	if htype{
		//before
		g.Tasks.onBefore = append(g.Tasks.onBefore,h)
	}else{
		//after
		g.Tasks.onAfter = append(g.Tasks.onAfter,h)
	}
	return g
}


func (g *GoRoutine)AddHandlersToChain(htypes []bool,hs []Handler)*GoRoutine{
	if len(htypes)!=len(hs){
		return g
	}else{
		g.lock.Lock()
		defer g.lock.Unlock()
		for k,v := range htypes{
			if v{
				//before
				g.Tasks.onBefore = append(g.Tasks.onBefore,hs[k])
			}else{
				//after
				g.Tasks.onAfter = append(g.Tasks.onAfter,hs[k])
			}
		}
		return g
	}
}

func (g *GoRoutine)Run()*GoRoutine{
	g.lock.Lock()
	defer g.lock.Unlock()
	if g.state ==GoRoutine_RunningStatus || g.state ==GoRoutine_WaitRunStatus || g.state == GoRoutine_IdleStatus{
		//已运行的不处理
		return g
	}

	if g.state == GoRoutine_EndStatus{
		//已终止的不能再次运行
		return g
	}

	runTimeChan := make(chan time.Duration)

	go func(){
		timeOut := time.NewTicker(g.maxAliveTime)
		defer timeOut.Stop()
ForEnd:
		for{
			//select 1
			select {
				case anyOp:=<-g.opChan:
					switch anyOp.opCode{
						case GoRoutine_OpCmd_End:
							//设置状态为终止状态
							g.state = GoRoutine_EndStatus
						case GoRoutine_OpCmd_Pause:
							g.state = GoRoutine_WaitRunStatus
						case GoRoutine_OpCmd_Resume:
							g.state = GoRoutine_IdleStatus
						default:

					}

				case anyTasks:=<-g.Tasks.queue:
					if g.state == GoRoutine_EndStatus{
						//如果状态已经是终止状态 则不处理新的任务  只把剩余的任务执行完
						continue
					}else if g.state == GoRoutine_WaitRunStatus {
						//如果是暂停状态 放到临时等待区
						g.tempTaskZone = append(g.tempTaskZone,anyTasks...)
						sort.Sort(g.tempTaskZone)
						continue
					}
					if len(anyTasks)>1 {
						//如果一次获取到多个任务  先按优先级排序
						sort.Sort(anyTasks)
						if g.state == GoRoutine_RunningStatus{
							//如果还是在执行状态 全放到临时区域
							g.tempTaskZone = append(g.tempTaskZone, anyTasks...)
							sort.Sort(g.tempTaskZone)
							continue
						}else if g.state == GoRoutine_IdleStatus {
							//如果空闲就取一个
							g.tempTaskZone = append(g.tempTaskZone, anyTasks[1:]...)
						}
					}

					//取出优先级最高的先执行
					task := anyTasks[0]
					execTask(g,task,runTimeChan)

				default:
					if (g.state == GoRoutine_IdleStatus || g.state == GoRoutine_EndStatus ) && len(g.tempTaskZone) >0 {
						//如果是空闲状态或者终止状态 且 临时区大小大于0
						task := g.tempTaskZone[0]
						g.tempTaskZone = g.tempTaskZone[1:]
						execTask(g,task,runTimeChan)
					}else if g.state == GoRoutine_EndStatus && len(g.tempTaskZone) ==0 {
						//所有任务都已经完成  可以退出
						close(runTimeChan)
						break ForEnd

					}
			}

			//select 2
			select {
				case <-timeOut.C:
					//超时 则把状态更改为终止状态 不接受新的任务  不强制终止
					g.state = GoRoutine_EndStatus
				default:
					if (g.state == GoRoutine_IdleStatus || g.state == GoRoutine_EndStatus ) && len(g.tempTaskZone) >0 {
						//如果是空闲状态或者终止状态 且 临时区大小大于0
						task := g.tempTaskZone[0]
						g.tempTaskZone = g.tempTaskZone[1:]
						execTask(g,task,runTimeChan)
					}else if g.state == GoRoutine_EndStatus && len(g.tempTaskZone) ==0 {
						//所有任务都已经完成  可以退出
						close(runTimeChan)
						break ForEnd

					}
			}
		}
	}()

	return g
}

func execHandlerChain(g *GoRoutine,htype bool)bool{

	param := make([]interface{},0)
	flag := false

	handleExecFunc := func(hs []Handler){
		for _,handle := range hs {
			flag,param = handle(param)
			if !flag{
				break
			}
		}
	}
	g.rlock.RLock()
	defer g.rlock.RUnlock()
	if htype{
		//before
		handleExecFunc(g.Tasks.onBefore)
	}else{
		//after
		handleExecFunc(g.Tasks.onAfter)
	}

	return flag
}


func execTask(g *GoRoutine,task *WorkTask,runTimeChan chan time.Duration){
	execHandlerChain(g,true)
	startTime :=time.Now()
	rs,err :=task.Tasks(task.Params)
	endTime := time.Now()
	usedTime := endTime.Sub(startTime)
	execHandlerChain(g,false)
	go func() {
		task.FutureResult <- MakeTaskResult(rs, err)
		close(task.FutureResult)
		runTimeChan <- usedTime
	}()
	g.historyJobs = append(g.historyJobs,task)
	//执行完成后更改状态 终止状态的情况 不能被更改
	if g.state !=GoRoutine_EndStatus {
		g.state = GoRoutine_IdleStatus
	}
}




/*
	分配一个工作任务
 */
func (g *GoRoutine)Dispatch(task *WorkTask){
	if g.state== GoRoutine_EndStatus{
		//如果是终止状态 则不再分配
		return
	}else {
		ts := make([]*WorkTask, 0)
		ts = append(ts, task)
		g.Tasks.queue <- ts
	}
}

/*
	分配一组工作任务
 */
func (g *GoRoutine)MultiDispatch(task []*WorkTask){
	if g.state== GoRoutine_EndStatus{
		//如果是终止状态 则不再分配
		return
	}else{
		g.Tasks.queue <- task
	}
}


