package common

import "testing"

func TestFormatCurrentDateYYYYMMdd(t *testing.T) {
	str := FormatCurrentDateYYYYMMdd()
	println(str)
}
