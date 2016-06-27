package Ldata

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"errors"
	"bytes"
)

type LateEngine struct {
	Ldb    *leveldb.DB
	DbName string
	DbPath string
}

func (ldb *LateEngine) Open() error {
	db, err := leveldb.OpenFile(ldb.DbPath, nil)
	if err != nil {
		return errors.New("Account database error")
	}
	//defer db.Close()
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
	ext, _ := ldb.Ldb.Has([]byte(key), nil)
	if ext {
		data, err := ldb.Ldb.Get([]byte(key), nil)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil,errors.New("not found")
}

func (ldb *LateEngine) Del(key string) error {
	ext, _ := ldb.Ldb.Has([]byte(key), nil)
	if ext{
		err := ldb.Ldb.Delete([]byte(key), nil)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("not found")
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