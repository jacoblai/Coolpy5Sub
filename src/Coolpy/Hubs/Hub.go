package Hubs

import (
	"github.com/garyburd/redigo/redis"
	"Coolpy/Incr"
	"encoding/json"
	"strconv"
	"strings"
	"errors"
	"Coolpy/Deller"
)

type Hub struct {
	Id        int64
	Ukey      string
	Title     string `validate:"required"`
	About     string
	Tags      []string
	Public    bool
	Local     string `validate:"required"`
	Latitude  float64 `validate:"gte=-90,lte=90"`
	Longitude float64 `validate:"gte=-180,lte=180"`
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
	rds.Do("SELECT", "2")
	go delChan()
}

func delChan() {
	for {
		select {
		case ukey, ok := <-Deller.DelHubs:
			if ok {
				hs, err := hubStartWith(ukey)
				if err != nil {
					break
				}
				for _, v := range hs {
					delhubs := ukey + ":" + strconv.FormatInt(v.Id, 10)
					del(delhubs)
					go deldos(delhubs)
				}
			}
		case delhub, ok := <-Deller.DelHub:
			if ok {
				del(delhub)
				go deldos(delhub)
			}
		}
		if Deller.DelHubs == nil && Deller.DelHub == nil {
			break
		}
	}
}

func deldos(ukeyhid string) {
	go func() {
		Deller.DelControls <- ukeyhid
	}()
	go func() {
		Deller.DelNodes <- ukeyhid
	}()
}

func hubCreate(hub *Hub) error {
	v, err := Incr.HubInrc()
	if err != nil {
		return err
	}
	hub.Id = v
	hub.Public = false
	json, err := json.Marshal(hub)
	if err != nil {
		return err
	}
	key := hub.Ukey + ":" + strconv.FormatInt(hub.Id, 10)
	_, err = rds.Do("SET", key, json)
	if err != nil {
		return err
	}
	return nil
}

func hubStartWith(k string) ([]*Hub, error) {
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
	_, err = rds.Do("SET", key, json)
	if err != nil {
		return err
	}
	return nil
}

func del(k string) error {
	if len(strings.TrimSpace(k)) == 0 {
		return errors.New("uid was nil")
	}
	_, err := redis.Int(rds.Do("DEL", k))
	if err != nil {
		return err
	}
	return nil
}

func All() ([]string, error) {
	data, err := redis.Strings(rds.Do("KEYS", "*"))
	if err != nil {
		return nil, err
	}
	return data, nil
}