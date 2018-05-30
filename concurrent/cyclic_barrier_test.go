package concurrent

import (
	"testing"
	"time"
	"fmt"
	"sync"
)

func TestNewCyclicBarrier(t *testing.T) {
	cbFunc := func()bool{
		fmt.Println("cbFunc exec ing ...")
		return true
	}
	i:=5
	cycBar := NewCyclicBarrier(i,&cbFunc)


	wg := sync.WaitGroup{}
	wg.Add(i)
	func1 := func(w *sync.WaitGroup){
		k := 7
		for k >0 {
			time.Sleep(time.Millisecond*10)
			fmt.Printf(">>>>>>>>>>>prepare ok ! %d \r\n",k)
			cycBar.Await()
			k--
		}
		wg.Done()
	}

	for i>0 {
		go func1(&wg)
		i--
	}
	wg.Wait()
	fmt.Println("<<<<<<<<<<<<<<<<<<<<")
}
