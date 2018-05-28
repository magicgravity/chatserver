package datastruct

import (
	"sync"
	"time"
	"errors"
)

const (
	Default_BlockingStack_PUSH_WaitTime = time.Second
	Default_BlockingStack_POP_WaitTime = time.Second
)


var (
	BlockingStack_NotInitialOk = errors.New("blockingStack have not initial ok")
	BlockingStack_OperTimeOut = errors.New("operate on blockingStack is time out")
	BlockingStack_OperFail = errors.New("operate on blockingStack is fail")
)

type BlockingStack struct {
	wlock *sync.Mutex
	rlock *sync.RWMutex
	size int
	topElement *interface{}
	bottomElement *interface{}
	container []interface{}
	pushWaitTimeLimit time.Duration
	popWaitTimeLimit time.Duration
	isInitalOk bool
}


func NewBlockingStack()*BlockingStack{
	bs := BlockingStack{}
	bs.size = 0
	bs.bottomElement = nil
	bs.topElement = nil
	bs.rlock  = &sync.RWMutex{}
	bs.wlock = &sync.Mutex{}
	bs.container = make([]interface{},0)
	bs.popWaitTimeLimit = Default_BlockingStack_POP_WaitTime
	bs.pushWaitTimeLimit = Default_BlockingStack_PUSH_WaitTime
	bs.isInitalOk = true
	return &bs
}


func (s *BlockingStack)Push(elm interface{})(bool,error){
	if s.isInitalOk {
		addOk := make(chan bool)
		timeOutFlag := false
		go func(flag *bool){
			s.wlock.Lock()
			s.rlock.Lock()
			if *flag {
				close(addOk)
			}else{
				s.container = append(s.container,elm)
				s.wlock.Unlock()
				s.size++
				s.topElement = &s.container[s.size-1]
				if s.size==1 {
					s.bottomElement = s.topElement
				}
				s.rlock.Unlock()
				addOk<-true
			}

		}(&timeOutFlag)
		select {
			case isOk:=<-addOk:
				if isOk {
					return true,nil
				}else{
					return false,BlockingStack_OperFail
				}
			case <-time.After(s.pushWaitTimeLimit):
				timeOutFlag = true
				return false,BlockingStack_OperTimeOut
		}

	}else{
		return false,BlockingStack_NotInitialOk
	}
}



func (s *BlockingStack)Pop()(bool,interface{},error){

	innerPop := func(s *BlockingStack)(bool,interface{},error){
		popOk := make(chan interface{})
		timeOutFlag := false
		go func(flag *bool){
			s.wlock.Lock()
			s.rlock.Lock()
			if *flag {
				close(popOk)
			}else{
				ret := s.container[s.size-1]
				s.container = s.container[0:s.size-1]
				s.wlock.Unlock()
				s.size--
				if s.size ==0 {
					s.topElement = nil
					s.bottomElement = nil
				}else{
					s.topElement = &s.container[s.size-1]
				}
				s.rlock.Unlock()
				popOk<-ret
			}
		}(&timeOutFlag)
		tC := time.NewTicker(s.popWaitTimeLimit)
		defer tC.Stop()
		select {
			case ret :=<-popOk:
				return true,ret,nil
			case <-tC.C:
				timeOutFlag = true
				return false,nil,BlockingStack_OperTimeOut
		}

	}

	if s.isInitalOk {
		s.rlock.RLock()
		if s.size ==0 {
			s.rlock.RUnlock()
			timeC := time.NewTicker(s.popWaitTimeLimit)
			defer timeC.Stop()
			for {
				select {
					case <-timeC.C:
						return false,nil,BlockingStack_OperTimeOut
					default:
						break
				}
				if s.size==0 {
					time.Sleep(time.Nanosecond*2)
					continue
				}else{
					break
				}
			}
			//空栈 变非空
			return innerPop(s)
		}else{
			s.rlock.RUnlock()
			return innerPop(s)
		}
	}else{
		return false,nil,BlockingStack_NotInitialOk
	}
}

func (s *BlockingStack)Size()int {
	s.rlock.RLock()
	defer s.rlock.RUnlock()
	return s.size
}


func (s *BlockingStack)SetTimeOut(pushTime,popTIme time.Duration)bool{
	if s.isInitalOk {
		s.popWaitTimeLimit = popTIme
		s.pushWaitTimeLimit = pushTime
		return true
	}else{
		return false
	}
}


func (s *BlockingStack)MultiPush(elems []interface{})(bool,error){
	if s.isInitalOk {
		addOk := make(chan bool)
		timeOutFlag := false
		go func(flag *bool){
			s.wlock.Lock()
			s.rlock.Lock()
			if *flag {
				close(addOk)
			}else{
				s.container = append(s.container,elems...)
				s.wlock.Unlock()
				s.size=s.size+len(elems)
				s.topElement = &s.container[s.size-1]
				if s.size==1 {
					s.bottomElement = s.topElement
				}
				s.rlock.Unlock()
				addOk<-true
			}

		}(&timeOutFlag)
		select {
		case isOk:=<-addOk:
			if isOk {
				return true,nil
			}else{
				return false,BlockingStack_OperFail
			}
		case <-time.After(s.pushWaitTimeLimit):
			timeOutFlag = true
			return false,BlockingStack_OperTimeOut
		}

	}else{
		return false,BlockingStack_NotInitialOk
	}
}


func (s *BlockingStack)Top()(bool,*interface{},error){
	if s.isInitalOk {
		s.rlock.RLock()
		defer s.rlock.RUnlock()
		if s.size>=1 {
			s.rlock.RLock()
			defer s.rlock.RUnlock()
			return true,s.topElement,nil
		}else {
			return false,nil,BlockingStack_OperFail
		}
	}else{
		return false,nil,BlockingStack_NotInitialOk
	}
}
