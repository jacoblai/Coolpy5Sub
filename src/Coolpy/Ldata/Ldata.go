package Ldata

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"errors"
	"bytes"
)

type LateEngine struct {
	Ldb    *leveldb.DB
	DbName string
	DbPath string
}

func (ldb *LateEngine) Open() error {
	o := &opt.Options{
		Compression:opt.NoCompression,
	}
	db, err := leveldb.OpenFile(ldb.DbPath, o)
	if err != nil {
		return errors.New("Account database error")
	}
	ldb.Ldb = db
	return nil
}

func (ldb *LateEngine) Set(key string, value []byte) error {
	err := ldb.Ldb.Put([]byte(key), value, nil)
	if err != nil {
		return err
	}
	return nil
}

func (ldb *LateEngine) Get(key string) ([]byte, error) {
	data, err := ldb.Ldb.Get([]byte(key), nil)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (ldb *LateEngine) Del(key string) error {
	err := ldb.Ldb.Delete([]byte(key), nil)
	if err != nil {
		return err
	}
	return nil
}

func (ldb *LateEngine) FindKeyStartWith(key string) (map[string][]byte, error) {
	keys := make(map[string][]byte)
	iter := ldb.Ldb.NewIterator(util.BytesPrefix([]byte(key)), nil)
	for iter.Next() {
		keys[string(iter.Key())] = iter.Value()
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return nil, errors.New("iter errors")
	} else {
		return keys, nil
	}
}

func (ldb *LateEngine) FindKeyRangeByDatetime(start string, end string) (map[string][]byte, error) {
	keys := make(map[string][]byte)
	min := []byte(start)
	max := []byte(end)
	iter := ldb.Ldb.NewIterator(nil, nil)
	for ok := iter.Seek(min); ok && bytes.Compare(iter.Key(), max) <= 0; ok = iter.Next() {
		keys[string(iter.Key())] = iter.Value()
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return nil, errors.New("iter errors")
	} else {
		return keys, nil
	}
}

func (ldb *LateEngine) FindAll() (map[string][]byte, error) {
	keys := make(map[string][]byte)
	iter := ldb.Ldb.NewIterator(nil, nil)
	for iter.Next() {
		keys[string(iter.Key())] = iter.Value()
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return nil, errors.New("iter errors")
	} else {
		return keys, nil
	}
}