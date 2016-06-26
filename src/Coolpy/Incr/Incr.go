package Incr

import (
	"Coolpy/LData"
	"encoding/binary"
	"sync"
)

var ldb *Ldata.LateEngine

func init() {
	db := &Ldata.LateEngine{DbPath:"data/incr", DbName:"IncrDb"}
	db.Open()
	//defer db.Ldb.Close()
	ldb = db
}

var bufPool = sync.Pool{
	New:func() interface{} {
		buf := make([]byte, 8)
		return buf
	},
}

func Incr(key string) uint64 {
	ext, _ := ldb.Ldb.Has([]byte(key), nil)
	buf := bufPool.Get().([]byte)
	defer bufPool.Put(buf)
	if !ext {
		binary.BigEndian.PutUint64(buf, 1)
		ldb.Ldb.Put([]byte(key), buf, nil)
		return 1
	} else {
		v, _ := ldb.Ldb.Get([]byte(key), nil)
		i := binary.BigEndian.Uint64(v) + 1
		binary.BigEndian.PutUint64(buf, i)
		ldb.Ldb.Put([]byte(key), buf, nil)
		return i
	}
}