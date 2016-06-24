package Account

import (
	//"github.com/golang/protobuf/proto"
	"github.com/satori/go.uuid"
	"encoding/json"
	"strings"
	"errors"
	"Coolpy/LData"
)

type Person struct {
	Ukey string
	Uid  string //required
	Pwd  string
}

var ldb *Ldata.LateEngine

func init() {
	db := &Ldata.LateEngine{DbPath:"data/ac", DbName:"AccountDb"}
	db.Open()
	//defer db.Ldb.Close()
	ldb = db
}

func New() *Person {
	return &Person{	}
}

func (p *Person) CreateUkey() {
	p.Ukey = uuid.NewV4().String()
}

func (p *Person) CreateOrReplace(ps *Person) error {
	if len(strings.TrimSpace(ps.Uid)) == 0 {
		return errors.New("uid was nil")
	}
	json, err := json.Marshal(ps)
	if err != nil {
		return err
	}
	if err = ldb.Set(ps.Uid, json); err != nil {
		return err
	}
	return nil
}

func (p *Person) Get(uid string) (*Person, error) {
	if len(strings.TrimSpace(uid)) == 0 {
		return nil, errors.New("uid was nil")
	}
	data, err := ldb.Get(uid)
	if err == nil {
		np := &Person{}
		json.Unmarshal(data, np)
		return np, nil
	}
	return nil, err
}

func (p *Person) Delete(uid string) error {
	if len(strings.TrimSpace(uid)) == 0 {
		return errors.New("uid was nil")
	}
	if err := ldb.Del(uid); err != nil {
		return err
	}
	return nil
}

func FindKeyStart(uid string) (map[string][]byte, error) {
	data, err := ldb.FindKeyStartWith(uid)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func All() (map[string][]byte, error) {
	data, err := ldb.FindAll()
	if err != nil {
		return nil, err
	}
	return data, nil
}