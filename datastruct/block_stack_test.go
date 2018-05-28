package datastruct

import (
	"testing"
	"fmt"
	"time"
	"sync"
)

func TestBlockingStack_Push(t *testing.T) {
	bs := NewBlockingStack()
	bs.Push(1)
	bs.Push(2)
	bs.Push(3)
	bs.Push(4)

	size := bs.Size()
	fmt.Printf("=====> %d \r\n",size)
	for ok,val,err := bs.Pop();;{
		if ok && err==nil {
			fmt.Printf("pop val =======>%v \r\n", val)
			ok,val,err = bs.Pop()
		}else{
			break
		}
	}
	size = bs.Size()
	fmt.Printf("=====> %d \r\n",size)
}



func TestBlockingStack_SetTimeOut(t *testing.T) {
	bs := NewBlockingStack()
	bs.Push(1)
	bs.Push(2)
	bs.Push(3)
	bs.Push(4)

	size := bs.Size()
	fmt.Printf("=====> %d \r\n",size)

	bs.SetTimeOut(time.Millisecond,time.Second*37)
	time1 := time.Now()
	for ok,val,err := bs.Pop();;{
		if ok && err==nil {
			fmt.Printf("pop val =======>%v \r\n", val)
			ok,val,err = bs.Pop()
		}else{
			break
		}
	}
	timeUsed := time.Now().Sub(time1)
	fmt.Printf("Used time >>> %v \r\n",timeUsed)
	size = bs.Size()
	fmt.Printf("=====> %d \r\n",size)
}


func TestBlockingStack_Push2(t *testing.T) {
	bs := NewBlockingStack()
	wg := sync.WaitGroup{}
	wg.Add(4)
	n := 10
	testF := func(a int){
		for i:=n;i>0;i--{
			bs.Push(i*a)
		}
		wg.Done()
	}

	go testF(2)
	go testF(3)
	go testF(4)
	go testF(5)
	wg.Wait()

	size := bs.Size()
	fmt.Printf("=====> %d \r\n",size)
	for ok,val,err := bs.Pop();;{
		if ok && err==nil {
			fmt.Printf("pop val =======>%v \r\n", val)
			ok,val,err = bs.Pop()
		}else{
			break
		}
	}
	size = bs.Size()
	fmt.Printf("=====> %d \r\n",size)
}