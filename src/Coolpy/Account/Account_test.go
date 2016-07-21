package Account

import (
	"testing"
	"fmt"
)

func TestAll(t *testing.T) {
	fmt.Println("start test")
	np := New()
	np.Uid = "jo111"
	np.Pwd = "pwd111afasdfasdfasdfasdfasdfasdfqfiweuriquweoruqowieruajsdfkajsdlkfjalskdjfklasjdfklajsdkl"
	np.CreateUkey()
	err := create(np)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println("create ok")
	}

	fmt.Println("get test")
	data, err := Get(np.Uid)
	if data.Uid == np.Uid {
		fmt.Println("find ok")
	} else {
		t.Error("find error")
	}

	fmt.Println("findkey test")
	fnp, _ := FindKeyStart("jo")
	if len(fnp) != 1 {
		t.Error("findkey error")
	}
	for _, v := range fnp {
		fmt.Println(v.Uid)
	}

	fmt.Println("findall test")
	anp, _ := all()
	if len(anp) != 1 {
		t.Error("allkey error")
	}
	for _, v := range anp {
		fmt.Println(v.Uid)
	}
}