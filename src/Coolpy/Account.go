package Coolpy

import (
	"github.com/satori/go.uuid"
	"strings"
	"errors"
	"regexp"
	"github.com/jacoblai/yiyidb"
)

type Person struct {
	Ukey     string `validate:"required"`
	Uid      string `validate:"required"`
	Pwd      string `validate:"required"`
	UserName string
	Email    string
}

func ValidateUidPwd(vstr string) bool {
	up := regexp.MustCompile("^[a-zA-Z0-9_]{3,128}$")
	re := up.MatchString(vstr)
	return re
}

var accrdsPool *yiyidb.Kvdb

func AccConnect(dir string) {
	db, err := yiyidb.OpenKvdb(dir+"/cp5accs", false, false, 10) //path, enable ttl
	if err != nil {
		panic(err)
	}
	accrdsPool = db
}

func AccNew() *Person {
	return &Person{}
}

func (p *Person) CreateUkey() {
	p.Ukey = uuid.NewV4().String()
}

func Acccreate(ps *Person) error {
	if len(strings.TrimSpace(ps.Uid)) == 0 {
		return errors.New("key nil")
	}
	return accrdsPool.PutJson([]byte(ps.Ukey+":"+ps.Uid), ps, 0)
}

func AccGet(uid string) (*Person, error) {
	if len(strings.TrimSpace(uid)) == 0 {
		return nil, errors.New("key nil")
	}
	k, err := AccgetFromDb(uid)
	if err != nil {
		return nil, err
	}

	var np Person
	err = accrdsPool.GetJson([]byte(k), &np)
	if err != nil {
		return nil, err
	}
	return &np, nil
}

func Accdel(uid string) error {
	if len(strings.TrimSpace(uid)) == 0 {
		return errors.New("key nil")
	}
	k, err := AccgetFromDb(uid)
	if err != nil {
		return err
	}
	return accrdsPool.Del([]byte(k))
}

func AccAll() ([]*Person) {
	var nt Person
	items := accrdsPool.AllByJson(nt)
	res := make([]*Person, 0)
	for _, v := range items {
		res = append(res, v.Object.(*Person))
	}
	return res
}

func AccgetFromDb(uid string) (string, error) {
	all := accrdsPool.AllKeys()
	for _, v := range all {
		if strings.HasSuffix(v, ":"+uid) {
			return string(v), nil
		}
	}
	return "", errors.New("not found")
}

func AccGetUkeyFromDb(k string) (string, error) {
	if len(strings.TrimSpace(k)) == 0 {
		return "", errors.New("ukey was nil")
	}
	items := accrdsPool.KeyStartKeys([]byte(k + ":"))
	if len(items) <= 0 {
		return "", errors.New("not ext")
	}
	return items[0], nil
}
