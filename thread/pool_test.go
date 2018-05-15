package thread

import (
	"testing"
	"fmt"
)

func TestGetExecutors(t *testing.T) {
	ms := make(map[int]*GoRoutine,2)
	fmt.Printf("======>%v",ms)
}
