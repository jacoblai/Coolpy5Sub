package Coolpy

import (
	"github.com/garyburd/redigo/redis"
	"encoding/json"
	"time"
	"github.com/pmylund/sortutil"
	"strings"
	"errors"
)

type GenDP struct {
	HubId     int64
	NodeId    int64
	TimeStamp time.Time
	Value     string `validate:"required"`
}

var GendprdsPool *redis.Pool

func GendpConnect(addr string, pwd string) {
	GendprdsPool = &redis.Pool{
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
			conn.Do("SELECT", "7")
			return conn, nil
		},
	}
}

func delGens(k string) {
	vs, err := GendpstartWith(k)
	if err != nil {
		return
	}
	for _, v := range vs {
		GendpDel(v)
	}
}

func GenCreate(k string, dp *GenDP) error {
	json, err := json.Marshal(dp)
	if err != nil {
		return err
	}
	rds := GendprdsPool.Get()
	defer rds.Close()
	_, err = rds.Do("SET", k, json)
	if err != nil {
		return err
	}
	return nil
}

func GendpstartWith(k string) ([]string, error) {
	rds := GendprdsPool.Get()
	defer rds.Close()
	data, err := redis.Strings(rds.Do("KEYSSTART", k))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GendpMaxGet(k string) (*GenDP, error) {
	rds := GendprdsPool.Get()
	defer rds.Close()
	data, err := redis.Strings(rds.Do("KEYSSTART", k))
	if err != nil {
		return nil, err
	}
	if len(data) <= 0 {
		return nil, errors.New("no data")
	}
	sortutil.Desc(data)
	o, _ := redis.String(rds.Do("GET", data[0]))
	dp := &GenDP{}
	err = json.Unmarshal([]byte(o), &dp)
	if err != nil {
		return nil, err
	}
	return dp, nil
}

func GendpGetOneByKey(k string) (*GenDP, error) {
	rds := GendprdsPool.Get()
	defer rds.Close()
	o, err := redis.String(rds.Do("GET", k))
	if err != nil {
		return nil, err
	}
	h := &GenDP{}
	err = json.Unmarshal([]byte(o), &h)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func GendpReplace(k string, h *GenDP) error {
	json, err := json.Marshal(h)
	if err != nil {
		return err
	}
	rds := GendprdsPool.Get()
	defer rds.Close()
	_, err = rds.Do("SET", k, json)
	if err != nil {
		return err
	}
	return nil
}

func GendpDel(k string) error {
	if len(strings.TrimSpace(k)) == 0 {
		return errors.New("uid was nil")
	}
	rds := GendprdsPool.Get()
	defer rds.Close()
	_, err := redis.Int(rds.Do("DEL", k))
	if err != nil {
		return err
	}
	return nil
}

func GendpGetRange(start string, end string, interval float64, page int) ([]*GenDP, error) {
	rds := GendprdsPool.Get()
	defer rds.Close()
	data, err := redis.Strings(rds.Do("KEYSRANGE", start, end))
	if err != nil {
		return nil, err
	}
	if len(data) <= 0 {
		return nil, errors.New("no data")
	}
	var IntervalData []string
	for _, v := range data {
		if len(IntervalData) == 0 {
			IntervalData = append(IntervalData, v)
		} else {
			otime := strings.Split(IntervalData[len(IntervalData) - 1], ",")
			otm, _ := time.Parse(time.RFC3339Nano, otime[2])
			vtime := strings.Split(v, ",")
			vtm, _ := time.Parse(time.RFC3339Nano, vtime[2])
			du := vtm.Sub(otm)
			if du.Seconds() >= interval {
				IntervalData = append(IntervalData, v)
			}
		}
	}
	var ndata []*GenDP
	for _, v := range IntervalData {
		o, _ := redis.String(rds.Do("GET", v))
		h := &GenDP{}
		json.Unmarshal([]byte(o), &h)
		ndata = append(ndata, h)
	}
	return ndata, nil
}

func GendpAll() ([]string, error) {
	rds := GendprdsPool.Get()
	defer rds.Close()
	data, err := redis.Strings(rds.Do("KEYS", "*"))
	if err != nil {
		return nil, err
	}
	return data, nil
}