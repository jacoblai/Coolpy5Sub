package Account

import (
	"github.com/satori/go.uuid"
	"encoding/json"
	"strings"
	"errors"
	"github.com/garyburd/redigo/redis"
	"regexp"
	"time"
)

type Person struct {
	Ukey     string `validate:"required"`
	Uid      string `validate:"required"`
	Pwd      string `validate:"required"`
	UserName string
	Email string
}

func ValidateUidPwd(vstr string) bool {
	up := regexp.MustCompile("^[a-zA-Z0-9_]{3,128}$")
	re := up.MatchString(vstr)
	return re
}

var rdsPool *redis.Pool

func Connect(addr string, pwd string) {
	rdsPool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: time.Second * 300,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			_, err = conn.Do("AUTH", pwd)
			if err != nil {
				return nil, err
			}
			conn.Do("SELECT", "1")
			return conn, nil
		},
	}
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
	rds := rdsPool.Get()
	defer rds.Close()
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
	rds := rdsPool.Get()
	defer rds.Close()
	o, err := redis.String(rds.Do("GET", k))
	if err != nil {
		return nil, err
	}
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
	rds := rdsPool.Get()
	defer rds.Close()
	_, err = redis.Int(rds.Do("DEL", k))
	if err != nil {
		return err
	}
	return nil
}

func All() ([]*Person, error) {
	rds := rdsPool.Get()
	defer rds.Close()
	data, err := redis.Strings(rds.Do("KEYS", "*"))
	if err != nil {
		return nil, err
	}
	if len(data) <= 0 {
		return nil, errors.New("no data")
	}
	var ndata []*Person
	for _, v := range data {
		if !strings.HasSuffix(v,"admin") {
			o, _ := redis.String(rds.Do("GET", v))
			np := &Person{}
			json.Unmarshal([]byte(o), &np)
			ndata = append(ndata, np)
		}
	}
	return ndata, nil
}

func getFromDb(uid string) (string, error) {
	rds := rdsPool.Get()
	defer rds.Close()
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
	rds := rdsPool.Get()
	defer rds.Close()
	data, err := redis.Strings(rds.Do("KEYS", k + ":*"))
	if err != nil {
		return "", err
	}
	if len(data) <= 0 {
		return "", errors.New("not ext")
	}
	return data[0], nil
}