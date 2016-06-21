package Account

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/satori/go.uuid"
	"fmt"
	"encoding/json"
)

type Person struct {
	Ukey string
	Uid  string
	Pwd  string
}

var ldb *leveldb.DB

func init() {
	db, err := leveldb.OpenFile("data/ac", nil)
	if err != nil {
		fmt.Println("Account database error:", err)
	}
	defer db.Close()
	ldb = db
}

func New() *Person {
	return &Person{
		Ukey : uuid.NewV4().String(),
	}
}

func (p *Person) Save(ps *Person) error {
	json, err := json.Marshal(ps)
	if err == nil {
		ldb.Put([]byte(ps.Uid), json, nil)
	}
	return err
}