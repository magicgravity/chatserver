package datastruct

import (
	"math/rand"
	"time"
)

type SkipList struct {
	level int
	top  *SkipListNode

}


//创建一个子节点
func createSkipListNode(element Comparable)SkipListNode{
	node := SkipListNode{}
	node.Element = element
	node.Next = nil
	node.DownNext = nil
	return node
}

//随机产生数k，k层下的都需要将值插入
func randomLevel(level int)int{
	k := 0
	rand.Seed(time.Now().Unix())
	for rand.Int()%2 == 0 {
		k++
	}
	if k>level{
		return level
	}else{
		return k
	}
}


func NewSkipList(level int)*SkipList{
	sl := SkipList{}
	sl.level =level
	var (
		tmpDown *SkipListNode = nil
		tmpNextDown *SkipListNode  =nil
		tmp *SkipListNode = nil
		tmpLevel = level
	)

	for ;tmpLevel != 0;tmpLevel--{
		node := createSkipListNode(MinComparableItem{})
		tmp = &node
		node = createSkipListNode(MaxComparableItem{})
		tmp.Next = &node
		tmp.DownNext = tmpDown
		tmp.Next.DownNext = tmpNextDown
		tmpDown = tmp
		tmpNextDown = tmp.Next
	}
	sl.top = tmp
	return &sl
}


func (sk SkipList)Find(target *Comparable) *SkipListNode{
	node := sk.top
	for {
		for node.Next.Element.compareTo(*target)<0 {
			node = node.Next
		}

		if node.DownNext == nil {
			return node
		}
		node = node.DownNext
	}
}


func (sk SkipList)Delete(target *Comparable)bool{
	tmpLevel := sk.level
	skNode := sk.top
	tmpNode := skNode
	flag := false
	for ;tmpLevel!=0 ; tmpLevel-- {
		for tmpNode.Next.Element.compareTo(*target) <0 {
			tmpNode  = tmpNode.Next
		}

		if tmpNode.Next.Element.compareTo(*target) ==0 {
			tmpNode.Next = tmpNode.Next.Next
			flag = true
		}

		tmpNode = skNode.DownNext
	}
	return flag
}


func (sk SkipList)Insert(target Comparable){
	var (
		skNode *SkipListNode = nil
		k = randomLevel(sk.level)
		tmp *SkipListNode = sk.top
		tmpLevel = sk.level
		tmpNode *SkipListNode = nil
		backTmpNode *SkipListNode = nil
		flag = 1
	)

	for ;tmpLevel!=k ; tmpLevel-- {
		tmp = tmp.DownNext
	}

	tmpLevel++
	tmpNode = tmp

	for ; tmpLevel!=0 ; tmpLevel--{
		for tmpNode.Next.Element.compareTo(target)<0 {
			tmpNode = tmpNode.Next
		}

		v := createSkipListNode(target)
		skNode = &v
		if flag!= -1 {
			backTmpNode.DownNext = skNode
		}
		backTmpNode = skNode
		skNode.Next = tmpNode.Next
		tmpNode.Next  = skNode
		flag = 0
		tmpNode = tmpNode.DownNext
	}

	return
}