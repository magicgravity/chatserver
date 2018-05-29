package concurrent

import (
	"testing"
	"fmt"
	"time"
	"sync"
)

func TestCountDownLatch_Await(t *testing.T) {
	totalC := 100
	begin := NewCountDownLatch(1)
	end := NewCountDownLatch(totalC)

	func1 := func(c int,beginLatch,endLatch *CountDownLatch){
		beginLatch.Await()
		fmt.Printf("========> %d \r\n",c)
		time.Sleep(time.Second)
		endLatch.CountDown()
	}


	for i:=0;i<totalC;i++{
		go func1(i,begin,end)
	}

	fmt.Println("开始倒数——————————————————>")
	time.Sleep(time.Second)
	begin.CountDown()
	end.Await()
	fmt.Println("倒计时结束----------------!")


}



func TestCountDownLatch_Await2(t *testing.T) {
	totalC := 10
	begin := sync.WaitGroup{}
	begin.Add(1)
	end := sync.WaitGroup{}
	end.Add(totalC)

	func1 := func(c int,beginLatch,endLatch *sync.WaitGroup){
		beginLatch.Wait()
		fmt.Printf("========> %d \r\n",c)
		time.Sleep(time.Second)
		endLatch.Done()
	}


	for i:=0;i<totalC;i++{
		go func1(i,&begin,&end)
	}

	fmt.Println("开始倒数——————————————————>")
	time.Sleep(time.Second)
	begin.Done()
	end.Wait()
	fmt.Println("倒计时结束----------------!")

}