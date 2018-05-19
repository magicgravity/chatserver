package common

import (
	"time"
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"unsafe"
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