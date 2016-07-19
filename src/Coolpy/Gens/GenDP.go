package Gens

import (
	"github.com/garyburd/redigo/redis"
	"encoding/json"
	"time"
	"github.com/pmylund/sortutil"
)

type GenDP struct {
	HubId     int64
	NodeId    int64
	TimeStamp time.Time
	Value     string `validate:"required"`
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
	rds.Do("SELECT", "7")
}

func GenCreate(k string, dp *GenDP) error {
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

func MaxGet(k string) (*GenDP, error) {
	data, err := redis.Strings(rds.Do("KEYSSTART", k))
	if err != nil {
		return nil, err
	}
	var ndata []*GenDP
	for _, v := range data {
		o, _ := redis.String(rds.Do("GET", v))
		h := &GenDP{}
		json.Unmarshal([]byte(o), &h)
		ndata = append(ndata, h)
	}
	sortutil.DescByField(ndata, "TimeStamp")
	return ndata[0], nil
}

func GetOneByKey(k string) (*GenDP, error) {
	o, err := redis.String(rds.Do("GET", k))
	if err != nil {
		return nil, err
	}
	h := &GenDP{}
	json.Unmarshal([]byte(o), &h)
	return h, nil
}