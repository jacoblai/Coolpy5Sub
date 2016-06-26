package Ldata

import (
	"testing"
	"fmt"
	"sort"
	"strings"
)

func TestLateEngine_FindKeyRangeByDatetime(t *testing.T) {
	db := &LateEngine{DbPath:"data/test", DbName:"AccountDb"}
	db.Open()
	//tm, err := time.Parse(time.RFC3339Nano, "2013-06-05T14:10:43.678Z")
	//if err != nil {
	//	panic(err)
	//}
	////t.Error(tm.Format(time.RFC3339Nano))
	//fmt.Println("datetime range test")
	//for i := 0; i < 10; i++ {
	//	key := tm.Add(time.Second * time.Duration(i))
	//	nkey := key.Format(time.RFC3339Nano)
	//	var nb []byte
	//	for _, r := range "1,2," {
	//		nb = append(nb, byte(r))
	//	}
	//	for _, r := range nkey {
	//		nb = append(nb, byte(r))
	//	}
	//	db.Ldb.Put(nb, []byte("0"), nil)
	//}

	//fmt.Println("datetime range read test")
	//all, _ := db.FindAll()
	//for k, _ := range all {
	//	fmt.Println(k)
	//}

	fmt.Println("datetime range read test")
	all, _ := db.FindKeyRangeByDatetime("1,2,2013-06-05T14:10:41", "1,2,2013-06-05T14:11:46")
	keys := make([]string, len(all))
	i := 0
	for k, _ := range all {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("Key:%v\n", strings.Split(k,",")[2])
	}
}
