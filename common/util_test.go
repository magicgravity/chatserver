package common

import (
	"testing"
	"fmt"
)

func TestFormatCurrentDateYYYYMMdd(t *testing.T) {
	str := FormatCurrentDateYYYYMMdd()
	println(str)
}


func TestGenMd5Result(t *testing.T) {
	md5str,err := GenMd5Result("浏览是的32323","")
	if err!=nil{
		t.Fatal(err)
	}else{
		fmt.Printf("======== %s",md5str)
	}
}

func TestGenRandomInt(t *testing.T) {
	v1 := GenRandomInt(3,12)
	v2 := GenRandomInt(10,3232)
	v3 := GenRandomInt(30,30)
	v4 := GenRandomInt(1000,32)
	v5 := GenRandomInt(99,99999999999)
	fmt.Printf("v1 == > %d \r\n",v1)
	fmt.Printf("v2 == > %d \r\n",v2)
	fmt.Printf("v3 == > %d \r\n",v3)
	fmt.Printf("v4 == > %d \r\n",v4)
	fmt.Printf("v5 == > %d \r\n",v5)
}

func TestSimpleStringMatch(t *testing.T) {
	var(
		testData = make(map[string]string)
		testExpectRes = make(map[string]bool)
	)
	testData["abcd1234453"]="cd"
	testExpectRes["abcd1234453"]=true

	testData["A2Bs2d"]="2Bs"
	testExpectRes["A2Bs2d"] = true


	testData["A2Bs2d3"]="2Bs232"
	testExpectRes["A2Bs2d3"] = false


	for k,v := range testData{
		fmt.Printf("srcStr======>%s ; patternStr====>%s  ",k,v)
		ret := SimpleStringMatch(v,k)
		fmt.Printf("======>%d  ",ret)
		exp,_ := testExpectRes[k]
		if ret>0 && exp {
			fmt.Println("ok")
		} else if ret <0 && !exp{
			fmt.Println("ok")
		}else{
			fmt.Println("fail")
		}
	}

}


func TestKMP(t *testing.T) {
	var(
		testData = make(map[string]string)
		testExpectRes = make(map[string]bool)
	)
	testData["abcd1234453"]="cd"
	testExpectRes["abcd1234453"]=true

	testData["A2Bs2d"]="2Bs"
	testExpectRes["A2Bs2d"] = true


	testData["A2Bs2d3"]="2Bs232"
	testExpectRes["A2Bs2d3"] = false


	for k,v := range testData{
		fmt.Printf("srcStr======>%s ; patternStr====>%s  ",k,v)
		ret := KMP(v,k)
		fmt.Printf("======>%d  ",ret)
		exp,_ := testExpectRes[k]
		if ret>0 && exp {
			fmt.Println("ok")
		} else if ret <0 && !exp{
			fmt.Println("ok")
		}else{
			fmt.Println("fail")
		}
	}

}

type IntArr []int

func (ia IntArr)Len() int{
	return len(ia)
}

func (ia IntArr)Less(i, j int) bool{
	return ia[i]<ia[j]
}

func (ia IntArr)Swap(i, j int){
	ia[i],ia[j] = ia[j],ia[i]
}

func (ia IntArr)Slice(s,e int)DataInterface{
	return ia[s:e]
}

func (ia IntArr)Equal(i,j int)bool{
	return ia[i]==ia[j]
}


func TestBasicBubbleSort(t *testing.T) {
	testData := []int{10,33,1,20031,232,999,11,8,23,65,34,-8,9}
	testArr := make(IntArr,0)
	testArr = append(testArr,testData...)

	BasicBubbleSort(testArr,true)

	for i,k := range testArr{
		fmt.Printf("[%d]----->[%d] \r\n",i,k)
	}
}

func TestAdvanceBubbleSortVer_01(t *testing.T) {
	testData := []int{10,33,1,20031,232,999,11,8,23}
	testArr := make(IntArr,0)
	testArr = append(testArr,testData...)

	AdvanceBubbleSortVer_01(testArr,true)

	for i,k := range testArr{
		fmt.Printf("[%d]----->[%d] \r\n",i,k)
	}
}


func TestAdvanceBubbleSortVer_02(t *testing.T) {
	testData := []int{10,33,1,20031,232,999,11,8,23}
	testArr := make(IntArr,0)
	testArr = append(testArr,testData...)

	AdvanceBubbleSortVer_02(testArr,true)

	for i,k := range testArr{
		fmt.Printf("[%d]----->[%d] \r\n",i,k)
	}
}

func TestAdvanceBubbleSortVer_03(t *testing.T) {
	testData := []int{10,33,1,20031,232,999,11,8,23}
	testArr := make(IntArr,0)
	testArr = append(testArr,testData...)

	AdvanceBubbleSortVer_03(testArr,false)

	for i,k := range testArr{
		fmt.Printf("[%d]----->[%d] \r\n",i,k)
	}
}


func TestChooseSort(t *testing.T) {
	testData := []int{10,33,1,20031,232,999,11,8,23}
	testArr := make(IntArr,0)
	testArr = append(testArr,testData...)

	ChooseSort(testArr,true)

	for i,k := range testArr{
		fmt.Printf("[%d]----->[%d] \r\n",i,k)
	}
}

func TestInsertSort(t *testing.T) {
	testData := []int{10,33,1,20031,232,999,11,8,23}
	testArr := make(IntArr,0)
	testArr = append(testArr,testData...)

	InsertSort(testArr,true)

	for i,k := range testArr{
		fmt.Printf("[%d]----->[%d] \r\n",i,k)
	}
}


func TestShellSort(t *testing.T) {
	testData := []int{10,33,1,20031,232,999,11,8,23}
	testArr := make(IntArr,0)
	testArr = append(testArr,testData...)

	ShellSort(testArr,true)

	for i,k := range testArr{
		fmt.Printf("[%d]----->[%d] \r\n",i,k)
	}
}


func TestMergeSort(t *testing.T) {
	testData := []int{10,33,1,20031,232,999,11,8,23,1023,1027,4,88}
	testArr := make(IntArr,0)
	testArr = append(testArr,testData...)

	MergeSort(testArr,true)

	for i,k := range testArr{
		fmt.Printf("[%d]----->[%d] \r\n",i,k)
	}
}

func TestQuickSort(t *testing.T) {
	testData := []int{10,33,1,20031,232,999,11,8,23,1023,1027,4,88}
	testArr := make(IntArr,0)
	testArr = append(testArr,testData...)

	QuickSort(testArr,false)

	for i,k := range testArr{
		fmt.Printf("[%d]----->[%d] \r\n",i,k)
	}

}


func TestGenRandomIntWithNegative(t *testing.T) {

	ranv := GenRandomIntWithNegative(-255, 30)
	fmt.Printf(">>>>>>>%d \r\n",ranv)

	ranv = GenRandomIntWithNegative(-15, 100)
	fmt.Printf(">>>>>>>%d \r\n",ranv)
}