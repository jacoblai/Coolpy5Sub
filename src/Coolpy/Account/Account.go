package Account

import (
	"github.com/satori/go.uuid"
	"encoding/json"
	"strings"
	"errors"
	"github.com/garyburd/redigo/redis"
)

type Person struct {
	Ukey     string `validate:"required"`
	Uid      string `validate:"required"`
	Pwd      string `validate:"required"`
	UserName string
}

var rds redis.Conn

func Connect(addr string, pwd string) {
	c, err := redis.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	_, err = c.Do("AUTH", pwd)
	if err != nil {
		panic(err)
	}
	rds = c
	_, err = redis.String(rds.Do("SELECT", "1"));
}

func New() *Person {
	return &Person{}
}

func (p *Person) CreateUkey() {
	p.Ukey = uuid.NewV4().String()
}

func create(ps *Person) error {
	if len(strings.TrimSpace(ps.Uid)) == 0 {
		return errors.New("uid was nil")
	}
	json, err := json.Marshal(ps)
	if err != nil {
		return err
	}
	k := ps.Ukey + ":" + ps.Uid
	_, err = rds.Do("SET", k, json)
	if err != nil {
		return err
	}
	return nil
}

func Get(uid string) (*Person, error) {
	if len(strings.TrimSpace(uid)) == 0 {
		return nil, errors.New("uid was nil")
	}
	k, err := getFromDb(uid)
	if err != nil {
		return nil, err
	}
	o, _ := redis.String(rds.Do("GET", k))
	np := &Person{}
	json.Unmarshal([]byte(o), &np)
	return np, nil
}

func GetByUkey(k string) (*Person, error) {
	if len(strings.TrimSpace(k)) == 0 {
		return nil, errors.New("uid was nil")
	}
	dbk, err := getUkeyFromDb(k)
	if err != nil {
		return nil, err
	}
	o, _ := redis.String(rds.Do("GET", dbk))
	np := &Person{}
	json.Unmarshal([]byte(o), &np)
	return np, nil
}

func del(uid string) error {
	if len(strings.TrimSpace(uid)) == 0 {
		return errors.New("uid was nil")
	}
	k, err := getFromDb(uid)
	if err != nil {
		return err
	}
	_, err = redis.Int(rds.Do("DEL", k))
	if err != nil {
		return err
	}
	return nil
}

func CheckUKey(k string) (bool, error) {
	data, err := redis.Strings(rds.Do("KEYSSTART", k))
	if err != nil {
		return false, err
	}
	if data[0] == "" {
		return false, nil
	}
	return true, nil
}

func FindKeyStart(uid string) (map[string]*Person, error) {
	data, err := redis.Strings(rds.Do("KEYSSTART", uid))
	if err != nil {
		return nil, err
	}
	ndata := make(map[string]*Person)
	for _, v := range data {
		o, _ := redis.String(rds.Do("GET", v))
		np := &Person{}
		json.Unmarshal([]byte(o), &np)
		ndata[v] = np
	}
	return ndata, nil
}

func all() ([]*Person, error) {
	data, err := redis.Strings(rds.Do("KEYS", "*"))
	if err != nil {
		return nil, err
	}
	var ndata []*Person
	for _, v := range data {
		o, _ := redis.String(rds.Do("GET", v))
		np := &Person{}
		json.Unmarshal([]byte(o), &np)
		ndata = append(ndata, np)
	}
	return ndata, nil
}

func getFromDb(uid string) (string, error) {
	data, err := redis.Strings(rds.Do("KEYS", "*:" + uid))
	if err != nil {
		return "", err
	}
	for _, v := range data {
		return v, nil
	}
	return "", errors.New("not ext")
}

func getUkeyFromDb(k string) (string, error) {
	data, err := redis.Strings(rds.Do("KEYS", k + ":*"))
	if err != nil {
		return "", err
	}
	for _, v := range data {
		return v, nil
	}
	return "", errors.New("not ext")
}