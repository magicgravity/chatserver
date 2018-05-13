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