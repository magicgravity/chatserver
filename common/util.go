package common

import (
	"time"
	"crypto/md5"
	"encoding/hex"
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