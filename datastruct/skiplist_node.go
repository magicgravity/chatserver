package datastruct


type Comparable interface {
	/*
	current compareTo another
	 1 ----> c > a
	 0 ----> c = a
	-1 ----> c < a
	 */
	compareTo(another Comparable)int
}

type SkipListNode struct {
	Element Comparable
	Next *SkipListNode
	DownNext *SkipListNode
}

type MinComparableItem struct {

}

type MaxComparableItem struct {

}

func (min MinComparableItem)compareTo(another Comparable)int{
	return -1
}

func (max MaxComparableItem)compareTo(another Comparable)int{
	return 1
}