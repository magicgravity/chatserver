package common

import (
	"time"
	"sync"
)

type Server interface {
	SetPort(string)
	GetPort()string
	Start()
}

type BaseServer struct {
	/*
		端口   ":50051"
	 */
	port string
	/*
	 是否启动 true-启动;false-停止
	 */
	isStart bool
	/*
	 启动时间
	 */
	startTime time.Time

	lock sync.Mutex
}

func (bs *BaseServer)SetPort(port string){
	bs.port = port
}

func (bs *BaseServer)GetPort()string{
	return bs.port
}

func (bs *BaseServer)Start()  {
	bs.lock.Lock()
	defer bs.lock.Unlock()
	if !bs.isStart {
		bs.isStart = true
		bs.startTime = time.Now()
	}
}

