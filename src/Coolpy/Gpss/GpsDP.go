package Gpss

import (
	"time"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
)

type GpsDP struct {
	HubId     int64
	NodeId    int64
	TimeStamp time.Time
	Lat float64 `validate:"required,gte=-90,lte=90"`
	Lng float64 `validate:"required,gte=-180,lte=180"`
	Speed int
	Offset int
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
	rds.Do("SELECT", "6")
}

func GpsCreate(k string, dp *GpsDP) error {
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