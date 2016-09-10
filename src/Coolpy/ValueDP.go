package Coolpy

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
	Value     float64
}

var valdprdsPool *redis.Pool

func ValdpConnect(addr string, pwd string) {
	valdprdsPool = &redis.Pool{
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
			conn.Do("SELECT", "5")
			return conn, nil
		},
	}
}

func delValues(k string) {
	vs, err := ValdpstartWith(k)
	if err != nil {
		return
	}
	for _, v := range vs {
		ValdpDel(v)
	}
}

func ValueCreate(k string, dp *ValueDP) error {
	json, err := json.Marshal(dp)
	if err != nil {
		return err
	}
	rds := valdprdsPool.Get()
	defer rds.Close()
	_, err = rds.Do("SET", k, json)
	if err != nil {
		return err
	}
	return nil
}

func ValdpstartWith(k string) ([]string, error) {
	rds := valdprdsPool.Get()
	defer rds.Close()
	data, err := redis.Strings(rds.Do("KEYSSTART", k))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func ValdpMaxGet(k string) (*ValueDP, error) {
	rds := valdprdsPool.Get()
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
	dp := &ValueDP{}
	err = json.Unmarshal([]byte(o), &dp)
	if err != nil {
		return nil, err
	}
	return dp, nil
}

func ValdpGetOneByKey(k string) (*ValueDP, error) {
	rds := valdprdsPool.Get()
	defer rds.Close()
	o, err := redis.String(rds.Do("GET", k))
	if err != nil {
		return nil, err
	}
	h := &ValueDP{}
	err = json.Unmarshal([]byte(o), &h)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func ValdpReplace(k string, h *ValueDP) error {
	json, err := json.Marshal(h)
	if err != nil {
		return err
	}
	rds := valdprdsPool.Get()
	defer rds.Close()
	_, err = rds.Do("SET", k, json)
	if err != nil {
		return err
	}
	return nil
}

func ValdpDel(k string) error {
	if len(strings.TrimSpace(k)) == 0 {
		return errors.New("uid was nil")
	}
	rds := valdprdsPool.Get()
	defer rds.Close()
	_, err := redis.Int(rds.Do("DEL", k))
	if err != nil {
		return err
	}
	return nil
}

func ValdpGetRange(start string, end string, interval float64, page int) ([]*ValueDP, error) {
	rds := valdprdsPool.Get()
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
	//pageSize := 50
	//allcount := len(IntervalData)
	//lastPageSize := allcount % pageSize
	//totalPage := (allcount + pageSize - 1) / pageSize
	//if page > totalPage {
	//	return nil, errors.New("pages out of range")
	//}
	//var pageData []string
	//if page == 1 {
	//	if totalPage == page {
	//		//只有一页
	//		pageData = IntervalData[:allcount]
	//	} else {
	//		//不止一页的第一页
	//		pageData = IntervalData[:pageSize]
	//	}
	//} else if page < totalPage {
	//	//中间页
	//	cursor := (pageSize * page) - pageSize //起启位计算
	//	pageData = IntervalData[cursor:cursor + pageSize]
	//} else if page == totalPage {
	//	//尾页
	//	if lastPageSize == 0 {
	//		pageData = IntervalData[allcount - pageSize:]
	//	} else {
	//		pageData = IntervalData[allcount - lastPageSize:]
	//	}
	//} else {
	//	return nil, errors.New("page not ext")
	//}
	var ndata []*ValueDP
	for _, v := range IntervalData {
		o, _ := redis.String(rds.Do("GET", v))
		h := &ValueDP{}
		json.Unmarshal([]byte(o), &h)
		ndata = append(ndata, h)
	}
	return ndata, nil
}

func ValdpAll() ([]string, error) {
	rds := valdprdsPool.Get()
	defer rds.Close()
	data, err := redis.Strings(rds.Do("KEYS", "*"))
	if err != nil {
		return nil, err
	}
	return data, nil
}