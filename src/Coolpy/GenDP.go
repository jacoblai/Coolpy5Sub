package Coolpy

import (
	"time"
	"github.com/pmylund/sortutil"
	"strings"
	"errors"
	"github.com/jacoblai/yiyidb"
)

type GenDP struct {
	HubId     int64
	NodeId    int64
	TimeStamp time.Time
	Value     string `validate:"required"`
}

var GendprdsPool *yiyidb.Kvdb

func GendpConnect(dir string) {
	db, err := yiyidb.OpenKvdb(dir+"/cp5gens", false, false, 10) //path, enable ttl
	if err != nil {
		panic(err)
	}
	GendprdsPool = db
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
	return GendprdsPool.PutJson([]byte(k), dp, 0)
}

func GendpstartWith(k string) ([]string, error) {
	ks := GendprdsPool.KeyStartKeys([]byte(k))
	return ks, nil
}

func GendpMaxGet(k string) (*GenDP, error) {
	ks := GendprdsPool.KeyStartKeys([]byte(k))
	if len(ks) <= 0 {
		return nil, errors.New("no data")
	}
	sortutil.Desc(ks)
	var dp GenDP
	err := GendprdsPool.GetJson([]byte(ks[0]), &dp)
	if err != nil {
		return nil, err
	}
	return &dp, nil
}

func GendpGetOneByKey(k string) (*GenDP, error) {
	var dp GenDP
	err := GendprdsPool.GetJson([]byte(k), &dp)
	if err != nil {
		return nil, err
	}
	return &dp, nil
}

func GendpReplace(k string, h *GenDP) error {
	return GendprdsPool.PutJson([]byte(k), h, 0)
}

func GendpDel(k string) error {
	if len(strings.TrimSpace(k)) == 0 {
		return errors.New("key nil")
	}
	return GendprdsPool.Del([]byte(k))
}

func GendpGetRange(start string, end string, interval float64, page int) ([]*GenDP, error) {
	var gdp GenDP
	data, err := GendprdsPool.KeyRangeByObject([]byte(start), []byte(end), gdp)
	if err != nil {
		return nil, err
	}
	if len(data) <= 0 {
		return nil, errors.New("no data")
	}
	var IntervalData []string
	for _, v := range data {
		if len(IntervalData) == 0 {
			IntervalData = append(IntervalData, string(v.Key))
		} else {
			otime := strings.Split(IntervalData[len(IntervalData)-1], ",")
			otm, _ := time.Parse(time.RFC3339Nano, otime[2])
			vtime := strings.Split(string(v.Key), ",")
			vtm, _ := time.Parse(time.RFC3339Nano, vtime[2])
			du := vtm.Sub(otm)
			if du.Seconds() >= interval {
				IntervalData = append(IntervalData, string(v.Key))
			}
		}
	}
	var ndata []*GenDP
	for _, v := range IntervalData {
		var h GenDP
		GendprdsPool.GetJson([]byte(v), &h)
		ndata = append(ndata, &h)
	}
	return ndata, nil
}

func GendpAll() ([]string, error) {
	ks := GendprdsPool.AllKeys()
	return ks, nil
}
