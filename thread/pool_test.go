package thread

import (
	"testing"
)

func TestGetExecutors(t *testing.T) {
	pe,err :=GetExecutors().newFixedThreadPool(100)
	if err != nil{
		t.Fatal(err)
	}else{
		pe.start()

	}
}