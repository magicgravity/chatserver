package datastruct

import (
	"sync"
	"fmt"
)

type Key interface {
	Hash()int
}

type MapNodeListNode struct {
	Key Key
	Data interface{}
	NextRef *MapNodeListNode
}

type mapSegment struct {
	innerMap map[int]MapNodeListNode
	rlock *sync.RWMutex
	wlock *sync.Mutex
}

const DefaultConcurrentHashMapSegmentSize  = 64

type ConcurrentHashMap struct {
	segmentSize int
	data []mapSegment
	totalItemCoount int
	rwlock *sync.RWMutex
	elementSize int
}

func (m *ConcurrentHashMap)calSegment(key Key)(int,int){
	hash := key.Hash()
	return hash%m.segmentSize,hash
}

func NewDefaultConcurrentHashMap()*ConcurrentHashMap{
	return NewConcurrentHashMap(DefaultConcurrentHashMapSegmentSize)
}

func NewConcurrentHashMap(size int)*ConcurrentHashMap{
	ccmap :=ConcurrentHashMap{}
	ccmap.rwlock = &sync.RWMutex{}
	ccmap.totalItemCoount = 0
	ccmap.data = make([]mapSegment,0)
	if size <DefaultConcurrentHashMapSegmentSize{
		ccmap.segmentSize = DefaultConcurrentHashMapSegmentSize
	}else{
		ccmap.segmentSize = size
	}
	ccmap.elementSize = 0
	for i:=0;i< ccmap.segmentSize;i++{
		ccmap.data = append(ccmap.data,newMapSegment())
	}
	return &ccmap
}

func newMapSegment()mapSegment{
	ms := mapSegment{}
	ms.wlock = &sync.Mutex{}
	ms.rlock = &sync.RWMutex{}
	ms.innerMap = make(map[int]MapNodeListNode)
	return ms
}


func incrElementSize(m *ConcurrentHashMap){
	m.rwlock.Lock()
	m.elementSize = m.elementSize+1
	m.rwlock.Unlock()
}

/*
	设置键值  如果已存在 则覆盖
 */
func (m *ConcurrentHashMap)Put(key Key,val interface{})bool{
	pos,hashcode := m.calSegment(key)
	fmt.Printf("pos==>[%d],hashcode==>[%d] \r\n",pos,hashcode)
	fmt.Printf("mdata size ==> %d \r\n",len(m.data))
	if pos <len(m.data) {
		m.rwlock.RLock()
		v := m.data[pos]
		m.rwlock.RUnlock()
		if len(v.innerMap)>0 {
			v.wlock.Lock()
			if vdata, ok := v.innerMap[hashcode]; ok {
				v.wlock.Unlock()
				if vdata.Key == key {
					//innerMap已经有该key  覆盖
					v.wlock.Lock()
					defer v.wlock.Unlock()
					vdata.Data = val
					incrElementSize(m)
					return false
				} else {
					var curNode  = &vdata
					for curNode.NextRef != nil {
						curNode := curNode.NextRef

						if curNode.Key == key{
							//innerMap 的list 已有该key 覆盖
							v.wlock.Lock()
							curNode.Data = val
							v.wlock.Unlock()
							incrElementSize(m)
							return false
						}
					}
					//innerMap 的list 没有该key  则在尾部增加新的节点
					newnode := MapNodeListNode{}
					newnode.NextRef = nil
					newnode.Key = key
					newnode.Data = val
					v.wlock.Lock()
					defer v.wlock.Unlock()
					curNode.NextRef = &newnode
					incrElementSize(m)
					return true

				}
			}else{
				v.wlock.Unlock()
				//innerMap还没有这个键 插入新的
				node := MapNodeListNode{}
				node.Key = key
				node.Data = val
				node.NextRef = nil
				v.wlock.Lock()
				defer v.wlock.Unlock()
				v.innerMap[hashcode] = node
				incrElementSize(m)
				return true
			}
		}else{
			//innerMap还没有这个键 插入新的
			node := MapNodeListNode{}
			node.Key = key
			node.Data = val
			node.NextRef = nil
			//if v.wlock==nil {
			//	v.wlock = &sync.Mutex{}
			//	v.rlock = &sync.RWMutex{}
			//	v.innerMap = make(map[int]MapNodeListNode)
			//}
			v.wlock.Lock()
			defer v.wlock.Unlock()
			v.innerMap[hashcode] = node
			incrElementSize(m)
			return true
		}
	}
	//还没有这个segment 新建一个
	m.rwlock.Lock()
	m.data[pos] = newMapSegment()
	m.rwlock.Unlock()
	m.data[pos].wlock.Lock()
	defer m.data[pos].wlock.Unlock()
	node := MapNodeListNode{}
	node.Key = key
	node.Data = val
	node.NextRef = nil
	m.data[pos].innerMap[hashcode] = node
	incrElementSize(m)
	return true

}



func (m *ConcurrentHashMap)Get(key Key)(bool,interface{}){
	pos,hashcode := m.calSegment(key)
	if pos <len(m.data) {
		//如果分区位置 在已存在的分区中
		m.rwlock.RLock()
		v := m.data[pos]
		m.rwlock.RUnlock()
		if len(v.innerMap)>0 {
			v.rlock.Lock()
			defer v.rlock.Unlock()
			if vdata, ok := v.innerMap[hashcode]; ok {
				if vdata.Key ==key{
					return true,vdata.Data
				}else{
					var curNode  = &vdata
					for curNode.NextRef != nil {
						curNode := curNode.NextRef

						if curNode.Key == key{
							return true,curNode.Data
						}
					}
					return false,nil
				}
			}
		}
	}
	return false,nil

}

func (m *ConcurrentHashMap)Size()int{
	m.rwlock.RLock()
	defer m.rwlock.RUnlock()
	return m.elementSize
}