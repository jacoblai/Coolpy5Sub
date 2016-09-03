package Coolpy

import (
	"github.com/garyburd/redigo/redis"
	"strconv"
	"encoding/json"
	"strings"
	"errors"
	"time"
)

type Switcher struct {
	HubId  int64
	NodeId int64
	Svalue int `validate:"required"`
}

type GenControl struct {
	HubId  int64
	NodeId int64
	Gvalue string `validate:"required"`
}

type RangeControl struct {
	HubId  int64
	NodeId int64
	Rvalue int64 `validate:"required"`
	Min    int64
	Max    int64
	Step   int64
}

var ctrlrdsPool *redis.Pool

func CtrlConnect(addr string, pwd string) {
	ctrlrdsPool = &redis.Pool{
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
			conn.Do("SELECT", "4")
			return conn, nil
		},
	}
}

func DelControls(k string) {
	cs, err := ctrlStartWith(k)
	if err != nil {
		return
	}
	for _, v := range cs {
		ctrldel(v)
	}
}

func ReplaceSwitcher(k string, s *Switcher) error {
	json, err := json.Marshal(s)
	if err != nil {
		return err
	}
	rds := ctrlrdsPool.Get()
	defer rds.Close()
	_, err = rds.Do("SET", k, json)
	if err != nil {
		return err
	}
	return nil
}

func GetSwitcher(k string) (*Switcher, error) {
	rds := ctrlrdsPool.Get()
	defer rds.Close()
	o, err := redis.String(rds.Do("GET", k))
	if err != nil {
		return nil, err
	}
	h := &Switcher{}
	err = json.Unmarshal([]byte(o), &h)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func BeginSwitcher(ukey string, Hubid int64, Nodeid int64) error {
	key := ukey + ":" + strconv.FormatInt(Hubid, 10) + ":" + strconv.FormatInt(Nodeid, 10)
	o := Switcher{
		HubId:Hubid,
		NodeId:Nodeid,
		Svalue:0,
	}
	json, err := json.Marshal(o)
	if err != nil {
		return err
	}
	rds := ctrlrdsPool.Get()
	defer rds.Close()
	_, err = rds.Do("SET", key, json)
	if err != nil {
		return err
	}
	return nil
}

func ReplaceRangeControl(k string, s *RangeControl) error {
	json, err := json.Marshal(s)
	if err != nil {
		return err
	}
	rds := ctrlrdsPool.Get()
	defer rds.Close()
	_, err = rds.Do("SET", k, json)
	if err != nil {
		return err
	}
	return nil
}

func GetRangeControl(k string) (*RangeControl, error) {
	rds := ctrlrdsPool.Get()
	defer rds.Close()
	o, err := redis.String(rds.Do("GET", k))
	if err != nil {
		return nil, err
	}
	h := &RangeControl{}
	err = json.Unmarshal([]byte(o), &h)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func BeginRangeControl(ukey string, Hubid int64, Nodeid int64) error {
	key := ukey + ":" + strconv.FormatInt(Hubid, 10) + ":" + strconv.FormatInt(Nodeid, 10)
	o := RangeControl{
		HubId:Hubid,
		NodeId:Nodeid,
		Rvalue:0,
		Min:0,
		Max:255,
		Step:5,
	}
	json, err := json.Marshal(o)
	if err != nil {
		return err
	}
	rds := ctrlrdsPool.Get()
	defer rds.Close()
	_, err = rds.Do("SET", key, json)
	if err != nil {
		return err
	}
	return nil
}

func ReplaceGenControl(k string, s *GenControl) error {
	json, err := json.Marshal(s)
	if err != nil {
		return err
	}
	rds := ctrlrdsPool.Get()
	defer rds.Close()
	_, err = rds.Do("SET", k, json)
	if err != nil {
		return err
	}
	return nil
}

func GetGenControl(k string) (*GenControl, error) {
	rds := ctrlrdsPool.Get()
	defer rds.Close()
	o, err := redis.String(rds.Do("GET", k))
	if err != nil {
		return nil, err
	}
	h := &GenControl{}
	err = json.Unmarshal([]byte(o), &h)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func BeginGenControl(ukey string, Hubid int64, Nodeid int64) error {
	key := ukey + ":" + strconv.FormatInt(Hubid, 10) + ":" + strconv.FormatInt(Nodeid, 10)
	o := GenControl{
		HubId:Hubid,
		NodeId:Nodeid,
		Gvalue:"",
	}
	json, err := json.Marshal(o)
	if err != nil {
		return err
	}
	rds := ctrlrdsPool.Get()
	defer rds.Close()
	_, err = rds.Do("SET", key, json)
	if err != nil {
		return err
	}
	return nil
}

func ctrlStartWith(k string) ([]string, error) {
	rds := ctrlrdsPool.Get()
	defer rds.Close()
	data, err := redis.Strings(rds.Do("KEYSSTART", k))
	if err != nil {
		return nil, err
	}
	if len(data) <= 0 {
		return nil, errors.New("no data")
	}
	return data, nil
}

func ctrldel(k string) error {
	if len(strings.TrimSpace(k)) == 0 {
		return errors.New("uid was nil")
	}
	rds := ctrlrdsPool.Get()
	defer rds.Close()
	_, err := redis.Int(rds.Do("DEL", k))
	if err != nil {
		return err
	}
	return nil
}

func CtrlAll() ([]string, error) {
	rds := ctrlrdsPool.Get()
	defer rds.Close()
	data, err := redis.Strings(rds.Do("KEYS", "*"))
	if err != nil {
		return nil, err
	}
	return data, nil
}