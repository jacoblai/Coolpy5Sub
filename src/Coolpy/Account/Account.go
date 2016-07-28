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
	Uid      string `validate:"required,min=3,max=18,regex=^[a-zA-Z0-9_]{3,18}$"`
	Pwd      string `validate:"required,min=3,max=18,regex=^[a-zA-Z0-9_]{3,18}$"`
	UserName string `validate:"regex=^[\\u4e00-\\u9fa5_a-zA-Z0-9-]{1,16}$"`
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
	err = json.Unmarshal([]byte(o), &np)
	if err != nil {
		return nil, err
	}
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

func All() ([]*Person, error) {
	data, err := redis.Strings(rds.Do("KEYS", "*"))
	if err != nil {
		return nil, err
	}
	if len(data) <= 0 {
		return nil, errors.New("no data")
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
	if len(data) <= 0 {
		return "", errors.New("not ext")
	}
	return data[0], nil
}

func GetUkeyFromDb(k string) (string, error) {
	if len(strings.TrimSpace(k)) == 0 {
		return "", errors.New("ukey was nil")
	}
	data, err := redis.Strings(rds.Do("KEYS", k + ":*"))
	if err != nil {
		return "", err
	}
	if len(data) <= 0 {
		return "", errors.New("not ext")
	}
	return data[0], nil
}