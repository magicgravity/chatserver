package concurrent

import (
	"testing"
	"fmt"
	"time"
)

func TestCountDownLatch_Await(t *testing.T) {
	totalC := 10
	begin := NewCountDownLatch(1)
	end := NewCountDownLatch(totalC)

	func1 := func(c int,beginLatch,endLatch *CountDownLatch){
		beginLatch.Await()
		fmt.Printf("========> %d \r\n",c)
		time.Sleep(time.Second*4)
		endLatch.CountDown()
	}


	for i:=0;i<totalC;i++{
		go func1(i,begin,end)
	}

	fmt.Println("开始倒数——————————————————>")
	begin.CountDown()
	end.Await()
	fmt.Println("倒计时结束----------------!")


}
