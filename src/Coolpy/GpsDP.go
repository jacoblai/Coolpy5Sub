package Coolpy

import (
	"time"
	"github.com/pmylund/sortutil"
	"strings"
	"errors"
	"github.com/jacoblai/yiyidb"
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

var gpsdprdsPool *yiyidb.Kvdb

func GpsdpConnect(dir string) {
	db, err := yiyidb.OpenKvdb(dir+"/cp5gpss", false, false, 10) //path, enable ttl
	if err != nil {
		panic(err)
	}
	gpsdprdsPool = db
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
	return gpsdprdsPool.PutJson([]byte(k), dp, 0)
}

func GpsdpstartWith(k string) ([]string, error) {
	ks := gpsdprdsPool.KeyStartKeys([]byte(k))
	return ks, nil
}

func GpsdpMaxGet(k string) (*GpsDP, error) {
	ks := gpsdprdsPool.KeyStartKeys([]byte(k))
	if len(ks) <= 0 {
		return nil, errors.New("no data")
	}
	sortutil.Desc(ks)
	var dp GpsDP
	err := gpsdprdsPool.GetJson([]byte(ks[0]), &dp)
	if err != nil {
		return nil, err
	}
	return &dp, nil
}

func GpsdpGetOneByKey(k string) (*GpsDP, error) {
	var dp GpsDP
	err := gpsdprdsPool.GetJson([]byte(k), &dp)
	if err != nil {
		return nil, err
	}
	return &dp, nil
}

func GpsdpReplace(k string, h *GpsDP) error {
	return gpsdprdsPool.PutJson([]byte(k), h, 0)
}

func GpsdpDel(k string) error {
	if len(strings.TrimSpace(k)) == 0 {
		return errors.New("key nil")
	}
	return gpsdprdsPool.Del([]byte(k))
}

func GpsdpGetRange(start string, end string, interval float64, page int) ([]*GpsDP, error) {
	var gdp GpsDP
	data, err := gpsdprdsPool.KeyRangeByObject([]byte(start), []byte(end), gdp)
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
	var ndata []*GpsDP
	for _, v := range IntervalData {
		var h GpsDP
		gpsdprdsPool.GetJson([]byte(v), &h)
		ndata = append(ndata, &h)
	}
	return ndata, nil
}

func GpsdpAll() ([]string, error) {
	ks := gpsdprdsPool.AllKeys()
	return ks, nil
}