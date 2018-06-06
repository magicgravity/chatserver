package main

import (
	"github.com/magicgravity/chatserver/thread"
	"log"
	"math"
)

func main(){

	pe,err :=thread.GetExecutors().NewFixedThreadPool(1000)
	if err!= nil {
		log.Fatal(err)
	}else{
		workNum := 1000
		testWork := func(params []interface{})interface{}{
			ilen := len(params)
			count :=0
			for i:=0;i<ilen;i++{
				val := params[i].(int)
				count += val*val
			}
			return count
		}
		paramSize := 1000000
		param := make([]interface{},paramSize)
		for i:=0;i<paramSize;i++{
			param[i] = math.Sqrt(float64(i*99))
		}
		resultSets := make([]chan *thread.TaskResult,workNum)
		for i:=0;i<workNum;i++ {
			chanRs,err := pe.Submit(1,testWork,param)
			if err!=nil{
				log.Fatal(err)
			}else{
				resultSets[i] = chanRs
			}
		}

		counter := 0
		for counter <workNum {
			for _, elm := range resultSets {
				select {
					case data := <-elm:
						log.Printf("=======>%v",data.String())
						counter++
				}
			}
		}


	}


}
