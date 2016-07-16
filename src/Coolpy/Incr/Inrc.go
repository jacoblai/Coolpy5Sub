package Incr

import (
	"github.com/garyburd/redigo/redis"
)

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
	rds.Do("SELECT", "0")
	if _, err := redis.String(c.Do("GET", "hubid")); err != nil {
		rds.Do("SET", "hubid", "0")
	}
}

func HubInrc() (int64, error) {
	v, err := redis.Int64(rds.Do("INCR", "hubid"))
	if err != nil {
		return 0, err
	}
	return v, nil
}