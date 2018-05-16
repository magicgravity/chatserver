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

