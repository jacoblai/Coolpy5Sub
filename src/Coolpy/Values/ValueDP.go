package Values

import (
	"github.com/garyburd/redigo/redis"
	"github.com/pmylund/sortutil"
	"encoding/json"
	"time"
	"strings"
	"errors"
)

type ValueDP struct {
	HubId     int64
	NodeId    int64
	TimeStamp time.Time
	Value     float64 `validate:"required"`
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
	rds.Do("SELECT", "5")
}

func ValueCreate(k string, dp *ValueDP) error {
	json, err := json.Marshal(dp)
	if err != nil {
		return err
	}
	_, err = rds.Do("SET", k, json)
	if err != nil {
		return err
	}
	return nil
}

func MaxGet(k string) (*ValueDP, error) {
	data, err := redis.Strings(rds.Do("KEYSSTART", k))
	if err != nil {
		return nil, err
	}
	var ndata []*ValueDP
	for _, v := range data {
		o, _ := redis.String(rds.Do("GET", v))
		h := &ValueDP{}
		json.Unmarshal([]byte(o), &h)
		ndata = append(ndata, h)
	}
	sortutil.DescByField(ndata, "TimeStamp")
	return ndata[0], nil
}

func GetOneByKey(k string) (*ValueDP, error) {
	o, err := redis.String(rds.Do("GET", k))
	if err != nil {
		return nil, err
	}
	h := &ValueDP{}
	json.Unmarshal([]byte(o), &h)
	return h, nil
}

func Replace(k string, h *ValueDP) error {
	json, err := json.Marshal(h)
	if err != nil {
		return err
	}
	_, err = rds.Do("SET", k, json)
	if err != nil {
		return err
	}
	return nil
}

func Delete(k string) error {
	if len(strings.TrimSpace(k)) == 0 {
		return errors.New("uid was nil")
	}
	_, err := redis.Int(rds.Do("DEL", k))
	if err != nil {
		return err
	}
	return nil
}