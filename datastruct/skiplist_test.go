package datastruct

import "testing"

type MyInt int

func (m MyInt)compareTo(another Comparable)int{
	if val,ok := another.(MyInt);ok {
		return int(m-val)
	}else{
		return 0
	}
}

func TestNewSkipList(t *testing.T) {
	list := NewSkipList(10)
	list.Insert(MyInt(1))
	//list.Insert(MyInt(2))
	//list.Insert(MyInt(5))
	//list.Insert(MyInt(1112))
	//list.Insert(MyInt(64))
	//list.Insert(MyInt(33))
	//list.Insert(MyInt(3))
	//list.Insert(MyInt(142))



}
