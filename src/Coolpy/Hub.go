package Coolpy

import (
	"github.com/garyburd/redigo/redis"
	"encoding/json"
	"strconv"
	"strings"
	"errors"
	"time"
	"fmt"
)

type Hub struct {
	Id        int64
	Ukey      string
	Title     string `validate:"required"`
	About     string
	Tags      []string
	Public    bool
	//Local     string `validate:"required"`
	//Latitude  float64 `validate:"gte=-90,lte=90"`
	//Longitude float64 `validate:"gte=-180,lte=180"`
}

var hubrdsPool *redis.Pool

func HubConnect(addr string, pwd string) {
	hubrdsPool = &redis.Pool{
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
			conn.Do("SELECT", "2")
			return conn, nil
		},
	}
}



func hubCreate(hub *Hub) error {
	v, err := HubInrc()
	if err != nil {
		return err
	}
	hub.Id = v
	hub.Public = false
	json, err := json.Marshal(hub)
	if err != nil {
		return err
	}
	rds := hubrdsPool.Get()
	defer rds.Close()
	key := hub.Ukey + ":" + strconv.FormatInt(hub.Id, 10)
	_, err = rds.Do("SET", key, json)
	if err != nil {
		return err
	}
	return nil
}

func hubStartWith(k string) ([]*Hub, error) {
	rds := hubrdsPool.Get()
	defer rds.Close()
	data, err := redis.Strings(rds.Do("KEYSSTART", k))
	if err != nil {
		return nil, err
	}
	if len(data) <= 0 {
		return nil, errors.New("no data")
	}
	var ndata []*Hub
	for _, v := range data {
		o, _ := redis.String(rds.Do("GET", v))
		h := &Hub{}
		json.Unmarshal([]byte(o), &h)
		h.Ukey = ""
		ndata = append(ndata, h)
	}
	return ndata, nil
}

func HubGetOne(k string) (*Hub, error) {
	rds := hubrdsPool.Get()
	defer rds.Close()
	o, err := redis.String(rds.Do("GET", k))
	if err != nil {
		return nil, err
	}
	h := &Hub{}
	err = json.Unmarshal([]byte(o), &h)
	if err != nil {
		return nil,err
	}
	return h, nil
}

func hubReplace(h *Hub) error {
	json, err := json.Marshal(h)
	if err != nil {
		return err
	}
	key := h.Ukey + ":" + strconv.FormatInt(h.Id, 10)
	rds := hubrdsPool.Get()
	defer rds.Close()
	_, err = rds.Do("SET", key, json)
	if err != nil {
		return err
	}
	return nil
}

func delhubs(ukey string) {
	hs, err := hubStartWith(ukey)
	if err != nil {
		fmt.Println("del hub err", err)
		return
	}
	for _, v := range hs {
		delhubs := ukey + ":" + strconv.FormatInt(v.Id, 10)
		hubdel(delhubs)
	}
}

func hubdel(k string) error {
	if len(strings.TrimSpace(k)) == 0 {
		return errors.New("uid was nil")
	}
	rds := hubrdsPool.Get()
	defer rds.Close()

	deldos(k)

	_, err := redis.Int(rds.Do("DEL", k))
	if err != nil {
		return err
	}
	return nil
}

func deldos(ukeyhid string) {
	DelControls(ukeyhid)
	delnodes(ukeyhid)
}

func hubAll() ([]string, error) {
	rds := hubrdsPool.Get()
	defer rds.Close()
	data, err := redis.Strings(rds.Do("KEYS", "*"))
	if err != nil {
		return nil, err
	}
	return data, nil
}