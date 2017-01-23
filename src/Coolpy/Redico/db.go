package Redico

import (
	"strconv"
	"sync"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"bytes"
	"github.com/syndtr/goleveldb/leveldb/filter"
)

const (
	defaultFilterBits int = 10
	KB int = 1024
	MB int = KB * 1024
	GB int = MB * 1024
)

// RedisDB holds a single (numbered) Redis database.
type RedicoDB struct {
	master       *sync.Mutex // pointer to the lock
	id           int         // db id
	leveldb      *leveldb.DB
	DbPath       string      //db path
	opts         *opt.Options
	iteratorOpts *opt.ReadOptions
	syncOpts     *opt.WriteOptions
}

func newRedicoDB(id int, l *sync.Mutex) RedicoDB {
	rdb := RedicoDB{
		id:            id,
		master:        l,
		DbPath: datapath + "/data/" + strconv.Itoa(id),
	}
	opts := &opt.Options{}
	opts.ErrorIfMissing = false
	opts.BlockCacheCapacity = 4 * MB
	opts.Filter = filter.NewBloomFilter(defaultFilterBits)
	opts.Compression = opt.NoCompression
	opts.BlockSize = 4 * KB
	opts.WriteBuffer = 4 * MB
	opts.OpenFilesCacheCapacity = 1 * KB
	opts.CompactionTableSize = 32 * MB
	opts.WriteL0SlowdownTrigger = 16
	opts.WriteL0PauseTrigger = 64
	rdb.opts = opts

	rdb.iteratorOpts = &opt.ReadOptions{}
	rdb.iteratorOpts.DontFillCache = true
	rdb.syncOpts = &opt.WriteOptions{}
	rdb.syncOpts.Sync = true

	ndb, err := leveldb.OpenFile(rdb.DbPath, rdb.opts)
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

//var bufPool = sync.Pool{
//	New:func() interface{} {
//		buf := make([]byte, 8)
//		return buf
//	},
//}

// change int key value
func (db *RedicoDB) stringIncr(k string, delta int) (int, error) {
	v := 0
	sv, err := db.leveldb.Get([]byte(k), nil)
	if err != nil {
		return 0, ErrKeyNotFound
	}
	v, err = strconv.Atoi(string(sv))
	if err != nil {
		return 0, ErrIntValueError
	}
	v += delta
	db.stringSet(k, strconv.Itoa(v))
	return v, nil
}

// allKeys returns all keys. Sorted.
func (db *RedicoDB) allKeys() []string {
	var keys []string
	iter := db.leveldb.NewIterator(nil, db.iteratorOpts)
	defer iter.Release()
	for iter.Next() {
		keys = append(keys, string(iter.Key()))
	}
	//sort.Strings(keys) // To make things deterministic.
	return keys
}

func (db *RedicoDB) keyStart(k string) []string {
	var keys []string
	iter := db.leveldb.NewIterator(util.BytesPrefix([]byte(k)), db.iteratorOpts)
	defer iter.Release()
	for iter.Next() {
		keys = append(keys, string(iter.Key()))
	}
	return keys
}

func (db *RedicoDB) keyRange(min string, max string) []string {
	var keys []string
	iter := db.leveldb.NewIterator(nil, db.iteratorOpts)
	defer iter.Release()
	for ok := iter.Seek([]byte(min)); ok && bytes.Compare(iter.Key(), []byte(max)) <= 0; ok = iter.Next() {
		keys = append(keys, string(iter.Key()))
	}
	return keys
}

func (db *RedicoDB) del(k string, delTTL bool) {
	if !db.exists(k) {
		return
	}
	err := db.leveldb.Delete([]byte(k), db.syncOpts)
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
	err := db.leveldb.Put([]byte(k), []byte(v), db.syncOpts)
	if err != nil {
		panic(err)
	}
}
