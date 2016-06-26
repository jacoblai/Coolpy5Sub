package Account

import (
	"testing"
	"fmt"
)

func TestAll(t *testing.T) {
	fmt.Println("start test")
	longnp := New()
	longnp.Uid = "li111"
	longnp.Pwd = "pwd111afasdfasdfasdfasdfasdfasdfqfiweuriquweoruqowieruajsdfkajsdlkfjalskdjfklasjdfklajsdkl"
	longnp.CreateUkey()
	err := CreateOrReplace(longnp)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println("create ok")
	}

	fmt.Println("get test")
	data, err := Get(longnp.Uid)
	if data.Uid == longnp.Uid {
		fmt.Println("find ok")
	} else {
		t.Error("find error")
	}

	fmt.Println("findkey test")
	fnp, _ := FindKeyStart("li")
	if len(fnp) != 1 {
		t.Error("findkey error")
	}
	for _, v := range fnp {
		fmt.Println(v.Uid)
	}

	fmt.Println("findall test")
	anp, _ := All()
	if len(anp) != 1{
		t.Error("allkey error")
	}
	for _, v := range anp {
		fmt.Println(v.Uid)
	}
}