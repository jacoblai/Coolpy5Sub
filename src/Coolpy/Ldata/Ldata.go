package Ldata

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
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
	data, err := ldb.Ldb.Get([]byte(key), nil)
	if err != nil{
		return nil,err
	}
	return data,nil
}