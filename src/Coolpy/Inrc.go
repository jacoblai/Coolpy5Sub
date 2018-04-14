package Coolpy

import (
	"github.com/jacoblai/yiyidb"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

var inrcrdsPool *yiyidb.Kvdb

func InrcConnect(dir string) {
	db, err := yiyidb.OpenKvdb(dir+"/cp5inrcs", false, false, 10) //path, enable ttl
	if err != nil {
		panic(err)
	}
	inrcrdsPool = db
	if _, err := inrcrdsPool.Get([]byte("hubid")); err != nil {
		inrcrdsPool.Put([]byte("hubid"), yiyidb.IdToKeyPure(0), 0)
	}
	if _, err := inrcrdsPool.Get([]byte("nodeid")); err != nil {
		inrcrdsPool.Put([]byte("nodeid"), yiyidb.IdToKeyPure(0), 0)
	}
}

func HubInrc() (uint64, error) {
	if o, err := inrcrdsPool.Get([]byte("hubid")); err == nil {
		oval := yiyidb.KeyToIDPure(o)
		oval++
		inrcrdsPool.Put([]byte("hubid"), yiyidb.IdToKeyPure(oval), 0)
		return oval, nil
	}
	return 0, errors.New("sys hub inrc err")
}

func NodeInrc() (uint64, error) {
	if o, err := inrcrdsPool.Get([]byte("nodeid")); err == nil {
		oval := yiyidb.KeyToIDPure(o)
		oval++
		inrcrdsPool.Put([]byte("nodeid"), yiyidb.IdToKeyPure(oval), 0)
		return oval, nil
	}
	return 0, errors.New("sys hub inrc err")
}
