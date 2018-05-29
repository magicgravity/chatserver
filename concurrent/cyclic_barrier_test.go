package concurrent

import (
	"testing"
	"time"
	"fmt"
)

func TestNewCyclicBarrier(t *testing.T) {
	cbFunc := func()bool{
		return true
	}
	cycBar := NewCyclicBarrier(2,&cbFunc)


	func1 := func(){
		time.Sleep(time.Second)
		fmt.Println(">>>>>>>>>>>prepare ok !")
		cycBar.Await()
	}
	i:=2
	for i>0 {
		go func1()
		i--
	}
	cycBar.Join()
	fmt.Println("<<<<<<<<<<<<<<<<<<<<")
}
