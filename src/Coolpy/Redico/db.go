package Redico

import (
	"sort"
	"strconv"
	"sync"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// RedisDB holds a single (numbered) Redis database.
type RedicoDB struct {
	master  *sync.Mutex // pointer to the lock
	id      int         // db id
	leveldb *leveldb.DB
	DbPath  string      //db path
}

func newRedicoDB(id int, l *sync.Mutex) RedicoDB {
	rdb := RedicoDB{
		id:            id,
		master:        l,
		DbPath: "data/" + strconv.Itoa(id),
	}
	o := &opt.Options{
		Compression:opt.NoCompression,
	}
	ndb, err := leveldb.OpenFile(rdb.DbPath, o)
	if err != nil {
		panic(err)
	}
	rdb.leveldb = ndb
	return rdb
}

func (db *RedicoDB) exists(k string) bool {
	ok, _ := db.leveldb.Has([]byte(k), nil)
	return ok
}

// allKeys returns all keys. Sorted.
func (db *RedicoDB) allKeys() []string {
	var keys []string
	iter := db.leveldb.NewIterator(nil, nil)
	for iter.Next() {
		keys = append(keys, string(iter.Key()))
	}
	iter.Release()
	sort.Strings(keys) // To make things deterministic.
	return keys
}

func (db *RedicoDB) del(k string, delTTL bool) {
	if !db.exists(k) {
		return
	}
	err := db.leveldb.Delete([]byte(k), nil)
	if err != nil {
		panic(err)
	}
}

// stringGet returns the string key or "" on error/nonexists.
func (db *RedicoDB) stringGet(k string) string {
	data, err := db.leveldb.Get([]byte(k), nil)
	if err != nil {
		return ""
	}
	return string(data)
}

// stringSet force set()s a key. Does not touch expire.
func (db *RedicoDB) stringSet(k, v string) {
	err := db.leveldb.Put([]byte(k), []byte(v), nil)
	if err != nil {
		panic(err)
	}
}
