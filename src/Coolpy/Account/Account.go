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

func Close() error {
	return rds.Close()
}

func (p *Person) CreateUkey() {
	p.Ukey = uuid.NewV4().String()
}

func createOrReplace(ps *Person) error {
	if len(strings.TrimSpace(ps.Uid)) == 0 {
		return errors.New("uid was nil")
	}
	json, err := json.Marshal(ps)
	if err != nil {
		return err
	}
	_, err = rds.Do("SET", ps.Uid, json)
	if err != nil {
		return err
	}
	return nil
}

func Get(uid string) (*Person, error) {
	if len(strings.TrimSpace(uid)) == 0 {
		return nil, errors.New("uid was nil")
	}
	data, err := redis.String(rds.Do("GET", uid))
	if err == nil {
		np := &Person{}
		json.Unmarshal([]byte(data), np)
		return np, nil
	}
	return nil, err
}

func delete(uid string) error {
	if len(strings.TrimSpace(uid)) == 0 {
		return errors.New("uid was nil")
	}
	_, err := redis.Int(rds.Do("DEL", uid))
	if err != nil {
		return err
	}
	return nil
}

func findKeyStart(uid string) (map[string]*Person, error) {
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