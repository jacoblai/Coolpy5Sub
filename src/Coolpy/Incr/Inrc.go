package Incr

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

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
			conn.Do("SELECT", "0")
			return conn, nil
		},
	}
	rds := rdsPool.Get()
	defer rds.Close()
	if _, err := redis.String(rds.Do("GET", "hubid")); err != nil {
		rds.Do("SET", "hubid", "0")
	}
	if _, err := redis.String(rds.Do("GET", "nodeid")); err != nil {
		rds.Do("SET", "nodeid", "0")
	}
}

func HubInrc() (int64, error) {
	rds := rdsPool.Get()
	defer rds.Close()
	v, err := redis.Int64(rds.Do("INCR", "hubid"))
	if err != nil {
		return 0, err
	}
	return v, nil
}

func NodeInrc() (int64, error) {
	rds := rdsPool.Get()
	defer rds.Close()
	v, err := redis.Int64(rds.Do("INCR", "nodeid"))
	if err != nil {
		return 0, err
	}
	return v, nil
}