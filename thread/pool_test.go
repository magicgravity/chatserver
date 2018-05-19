package thread

import (
	"testing"
	"fmt"
)

func TestGetExecutors(t *testing.T) {
	pe,err :=GetExecutors().newFixedThreadPool(100)
	if err != nil{
		t.Fatal(err)
	}else{
		pe.start()

		job1:= func(params []interface{})interface{}{
			ilen := len(params)
			count :=0
			for i:=0;i<ilen;i++{
				val := params[i].(int)
				count += val*val*val*count
			}
			return count
		}
		param := make([]interface{},100)
		for i:=0;i<100;i++{
			param[i] = i*i
		}
		chanRs,err := pe.submit(1,job1,param)
		if err!=nil{
			t.Fatal(err)
		}
		select{
			case taskRs:= <-chanRs:
				fmt.Printf("-------->%v",taskRs)
		}

	}
}