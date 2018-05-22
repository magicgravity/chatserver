package common

import (
	"time"
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"unsafe"
	"sort"
	"fmt"
)

const(
	c1_32 uint32 = 0xcc9e2d51
	c2_32 uint32 = 0x1b873593
)

func FormatCurrentDateYYYYMMdd()string{
	now := time.Now()
	return now.Format("20060102150405")
}


/*
计算Md5
 */
func GenMd5Result(raw ,salt string) (string,error){
	Md5 := md5.New()
	_,err :=Md5.Write(([]byte)(raw+salt))
	if err!=nil{
		return "",err
	}else{
		cipherStr :=Md5.Sum(nil)
		return hex.EncodeToString(cipherStr),nil
	}
}

/*
	生成一个ID序号
 */
func GenGoroutineId()int{
	//todo
	return 0
}


func GenRandomInt(min ,max int)int {
	rand.Seed(time.Now().Unix())
	if max>min {
		return 	min + rand.Intn(max - min)
	}else  if max == min {
		return  max
	}else {
		return max + rand.Intn(min-max)
	}
}



// GetHash returns a murmur32 hash for the data slice.
func GetHash(data []byte) uint32 {
	// Seed is set to 37, same as C# version of emitter
	var h1 uint32 = 37

	nblocks := len(data) / 4
	var p uintptr
	if len(data) > 0 {
		p = uintptr(unsafe.Pointer(&data[0]))
	}

	p1 := p + uintptr(4*nblocks)
	for ; p < p1; p += 4 {
		k1 := *(*uint32)(unsafe.Pointer(p))

		k1 *= c1_32
		k1 = (k1 << 15) | (k1 >> 17) // rotl32(k1, 15)
		k1 *= c2_32

		h1 ^= k1
		h1 = (h1 << 13) | (h1 >> 19) // rotl32(h1, 13)
		h1 = h1*5 + 0xe6546b64
	}

	tail := data[nblocks*4:]

	var k1 uint32
	switch len(tail) & 3 {
	case 3:
		k1 ^= uint32(tail[2]) << 16
		fallthrough
	case 2:
		k1 ^= uint32(tail[1]) << 8
		fallthrough
	case 1:
		k1 ^= uint32(tail[0])
		k1 *= c1_32
		k1 = (k1 << 15) | (k1 >> 17) // rotl32(k1, 15)
		k1 *= c2_32
		h1 ^= k1
	}

	h1 ^= uint32(len(data))

	h1 ^= h1 >> 16
	h1 *= 0x85ebca6b
	h1 ^= h1 >> 13
	h1 *= 0xc2b2ae35
	h1 ^= h1 >> 16

	return (h1 << 24) | (((h1 >> 8) << 16) & 0xFF0000) | (((h1 >> 16) << 8) & 0xFF00) | (h1 >> 24)
}


func SimpleStringMatch(patternstr,srcstr string)int{
	var srcStrLen = len(srcstr)
	var patternStrLen = len(patternstr)
	for i:=0;i<= srcStrLen-patternStrLen;i++{
		var j = 0
		for{
			if j<patternStrLen && patternstr[j]==srcstr[i+j]{
				j ++
			}else{
				break
			}
		}
		if j==patternStrLen{
			return i
		}
	}
	return -1
}


func KMP(patternStr ,srcStr string)int{
	var (
		patternLen = len(patternStr)
		srcLen = len(srcStr)
		prefixArr = make([]int,patternLen)
		i =0
		j =0
	)

	makePrefixTableFunc := func(p *string,plen int){
		var (
			i =1
			j =0
		)
		prefixArr[0] = 0
		for {
			if i<plen{
				if (*p)[i]==(*p)[j]{
					prefixArr[i]=j+1
					i++
					j++
				}else if j>0{
					j=prefixArr[j-1]
				}else{
					prefixArr[i]=0
					i++
				}
			}else{
				break
			}
		}
	}
	//
	makePrefixTableFunc(&patternStr,patternLen)

	for {
		if i< srcLen{
			if srcStr[i]==patternStr[j]{
				if j==patternLen-1	{
					return i-j
				}else{
					i++
					j++
				}
			}else if j>0 {
				j = prefixArr[j-1]
			}else{
				i++
			}
		}else{
			break
		}
	}

	return -1

}


/*
  基本的冒泡排序
 */
func BasicBubbleSort(arr sort.Interface,ascOrder bool)sort.Interface{
	arrLen := arr.Len()
	for i:=arrLen-1;i>=0;i-- {
		for j:=0;j<i;j++ {
			if ascOrder {
				if arr.Less(i, j) {
					arr.Swap(i, j)
				}
			}else{
				if arr.Less(j, i) {
					arr.Swap(i, j)
				}
			}
		}
	}
	return arr
}


/*
	改进后的冒泡排序
	增加了一个标记 可以通过标记判断提前结束
	todo ?????有问题
 */
func AdvanceBubbleSortVer_01(arr sort.Interface,ascOrder bool)sort.Interface{
	arrLen := arr.Len()
	isSorted := false
	for i:=arrLen-1;i>0 && !isSorted ;i-- {
		fmt.Printf("i:%d =====>%v \r\n",i,arr)
		isSorted = true
		for j :=0;j<i-1;j++ {
			fmt.Printf("i:%d   j:%d   ~~~~~~~~~>%v \r\n",i,j,arr)
			if ascOrder{
				if arr.Less(i,j){
					arr.Swap(i,j)
					isSorted = false
				}
			}else{
				if arr.Less(j,i){
					arr.Swap(i,j)
					isSorted = false
				}
			}
		}
	}
	fmt.Println()
	return arr

}


/*
	 每次循环记录最后一次发生交换的元素的位置，这说明这之后的元素已经有序，下一次循环不用比较这些元素。
     最好情况时间复杂度为O(n)，最坏和平均情况时间复杂度为O(n^2)。
 */
func AdvanceBubbleSortVer_02(arr sort.Interface,ascOrder bool)sort.Interface{
	last := arr.Len()-1
	cur := 0
	for last >0{
		cur = 0
		for i:=0;i<last;i++ {
			fmt.Printf("i:%d   ;   last:%d=====>%v \r\n",i,last,arr)
			if !ascOrder {
				if arr.Less(i, i+1) {
					arr.Swap(i, i+1)
					cur = i
				}
			}else{
				if arr.Less(i+1,i) {
					arr.Swap(i, i+1)
					cur = i
				}
			}
		}
		last = cur
	}
	return arr
}

/*
	双向扫描的冒泡排序(鸡尾酒排序)
	每次循环不仅从前向后扫描记录最后一次发生交换的元素的位置up，而且从后向前扫描记录再次扫描记录最前面发生交换的元素的位置low，
	这样两侧的元素已经有序，当low>=up的时候证明整个数组有序。
    最好情况时间复杂度为O(n)，最坏和平均情况时间复杂度为O(n^2)。
 */
func AdvanceBubbleSortVer_03(arr sort.Interface,ascOrder bool)sort.Interface{
	var(
		low = 0
		up = arr.Len()-1
		index = 0
		i = 0
	)
	for up>low{
		for i=low;i<up;i++ {
			if ascOrder {
				if arr.Less(i+1,i) {
					arr.Swap(i,i+1)
					index = i
				}
			}else{
				if arr.Less(i,i+1) {
					arr.Swap(i,i+1)
					index = i
				}
			}
		}
		up = index
		for i=up;i>low;i-- {
			if ascOrder {
				if arr.Less(i, i-1) {
					arr.Swap(i, i-1)
					index = i
				}
			}else{
				if arr.Less(i-1, i) {
					arr.Swap(i, i-1)
					index = i
				}
			}
		}
		low = index
	}

	return arr
}

/*
选择排序
 */
func ChooseSort(arr sort.Interface,ascOrder bool)sort.Interface{
	var (
		min = 0
		arrLen = arr.Len()
	)
	for i:=0;i<arrLen-1;i++ {
		min = i
		for j:=i+1;j<arrLen;j++ {
			if ascOrder {
				if arr.Less(j, min) {
					min = j
				}
			}else{
				if arr.Less(min,j) {
					min = j
				}
			}
		}
		arr.Swap(i,min)

	}
	return arr
}


/*
	插入排序
 */
func InsertSort(arr sort.Interface,ascOrder bool)sort.Interface{
	arrLen := arr.Len()

	for i:=1;i<arrLen;i++ {
		for j:=i-1;j>=0 && ((arr.Less(j+1,j)&&ascOrder) || (!ascOrder && arr.Less(j,j+1)));j-- {
			arr.Swap(j,j+1)
		}
	}

	return arr
}


/*
	希尔排序
 */
func ShellSort(arr sort.Interface,ascOrder bool)sort.Interface{
	var(
		arrLen = arr.Len()
		gap = arrLen/2
	)

	for gap>=1 {
		for i:=gap;i<arrLen;i=i+gap {
			for j:=i-gap;j>=0 && ((arr.Less(j+gap,j)&&ascOrder) || (!ascOrder && arr.Less(j,j+gap)));j=j-gap {
				arr.Swap(j,j+gap)
			}
		}
		gap = gap /2
	}
	return arr
}