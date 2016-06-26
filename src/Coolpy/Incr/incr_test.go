package Incr

import (
	"testing"
	"fmt"
)

func TestIncr(t *testing.T) {
	id := Incr("hubid")
	//if id !=2{
	//	t.Error("id init error")
	//}
	fmt.Println("hubid:",id)
	id = Incr("modeid")
	fmt.Println("nodeid:",id)
}
