package thread

import (
	"testing"
	"fmt"
)

func TestGetExecutors(t *testing.T) {
	pe,err :=GetExecutors().NewFixedThreadPool(100)
	if err != nil{
		t.Fatal(err)
	}else{
		pe.Start()

		job1:= func(params []interface{})interface{}{
			ilen := len(params)
			count :=0
			for i:=0;i<ilen;i++{
				val := params[i].(int)
				count += val*val
			}
			return count
		}
		param := make([]interface{},100)
		for i:=0;i<100;i++{
			param[i] = i*i
		}
		chanRs,err := pe.Submit(1,job1,param)
		if err!=nil{
			t.Fatal(err)
		}
		select{
			case taskRs:= <-chanRs:
				fmt.Printf("-------->%v",taskRs)
		}

	}
}