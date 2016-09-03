package Coolpy

import (
	"time"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/pmylund/sortutil"
	"strings"
	"errors"
)

type GpsDP struct {
	HubId     int64
	NodeId    int64
	TimeStamp time.Time
	Lat       float64 `validate:"required,gte=-90,lte=90"`
	Lng       float64 `validate:"required,gte=-180,lte=180"`
	Speed     int
	Offset    int
}

var gpsdprdsPool *redis.Pool

func GpsdpConnect(addr string, pwd string) {
	gpsdprdsPool = &redis.Pool{
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
			conn.Do("SELECT", "6")
			return conn, nil
		},
	}
}

func delGpss(k string) {
	vs, err := GpsdpstartWith(k)
	if err != nil {
		return
	}
	for _, v := range vs {
		GpsdpDel(v)
	}
}

func GpsCreate(k string, dp *GpsDP) error {
	json, err := json.Marshal(dp)
	if err != nil {
		return err
	}
	rds := gpsdprdsPool.Get()
	defer rds.Close()
	_, err = rds.Do("SET", k, json)
	if err != nil {
		return err
	}
	return nil
}

func GpsdpstartWith(k string) ([]string, error) {
	rds := gpsdprdsPool.Get()
	defer rds.Close()
	data, err := redis.Strings(rds.Do("KEYSSTART", k))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GpsdpMaxGet(k string) (*GpsDP, error) {
	rds := gpsdprdsPool.Get()
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
	dp := &GpsDP{}
	err = json.Unmarshal([]byte(o), &dp)
	if err != nil {
		return nil, err
	}
	return dp, nil
}

func GpsdpGetOneByKey(k string) (*GpsDP, error) {
	rds := gpsdprdsPool.Get()
	defer rds.Close()
	o, err := redis.String(rds.Do("GET", k))
	if err != nil {
		return nil, err
	}
	h := &GpsDP{}
	err = json.Unmarshal([]byte(o), &h)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func GpsdpReplace(k string, h *GpsDP) error {
	json, err := json.Marshal(h)
	if err != nil {
		return err
	}
	rds := gpsdprdsPool.Get()
	defer rds.Close()
	_, err = rds.Do("SET", k, json)
	if err != nil {
		return err
	}
	return nil
}

func GpsdpDel(k string) error {
	if len(strings.TrimSpace(k)) == 0 {
		return errors.New("uid was nil")
	}
	rds := gpsdprdsPool.Get()
	defer rds.Close()
	_, err := redis.Int(rds.Do("DEL", k))
	if err != nil {
		return err
	}
	return nil
}

func GpsdpGetRange(start string, end string, interval float64, page int) ([]*GpsDP, error) {
	rds := gpsdprdsPool.Get()
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
	var ndata []*GpsDP
	for _, v := range IntervalData {
		o, _ := redis.String(rds.Do("GET", v))
		h := &GpsDP{}
		json.Unmarshal([]byte(o), &h)
		ndata = append(ndata, h)
	}
	return ndata, nil
}

func GpsdpAll() ([]string, error) {
	rds := gpsdprdsPool.Get()
	defer rds.Close()
	data, err := redis.Strings(rds.Do("KEYS", "*"))
	if err != nil {
		return nil, err
	}
	return data, nil
}