package datastruct

import (
	"testing"
	"github.com/magicgravity/chatserver/common"
	"strconv"
	"fmt"
	"sync"
)

type CString string

func (ks CString)Hash()int{
	hashval := common.GetHash(([]byte)(string(ks)))
	return int(hashval)
}

func genData(size int)map[CString]interface{}{
	var i =0
	rst := make(map[CString]interface{},0)
	for{
		if i>=size{
			break
		}else{
			i++
			key := "key—"+strconv.Itoa(i+1)+strconv.Itoa(common.GenRandomInt(100,999))
			val := "value:"+strconv.Itoa(common.GenRandomInt(i,999*i))
			rst[CString(key)] = val
		}
	}
	return rst
}


func TestConcurrentHashMap_Put(t *testing.T) {
	cmap := NewConcurrentHashMap(10)
	wg := sync.WaitGroup{}
	var concurrentNum =2
	wg.Add(concurrentNum)
	testFunc := func() {
		dd := genData(60)
		for k, v := range dd {
			cmap.Put(k, v)
		}

		fmt.Printf("map size == [%d] \r\n", cmap.Size())

		for k, v := range dd {
			ok, vv := cmap.Get(k)
			if ok && vv == v {
				fmt.Print("find match ok ~~~~\r\n")
			} else {
				t.Fatalf("find match fail ! key==>[%s],val===>[%s],mapval===>[%s] \r\n", k, v, vv)
			}
		}
		wg.Done()
	}
	for i:=0;i<concurrentNum;i++ {
		go testFunc()
	}

	wg.Wait()
}
