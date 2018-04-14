package Coolpy

import (
	"time"
	"github.com/pmylund/sortutil"
	"strings"
	"errors"
	"github.com/jacoblai/yiyidb"
)

type PhotoDP struct {
	HubId     uint64
	NodeId    uint64
	TimeStamp time.Time
	Size      int64 `validate:"required"`
	Mime      string `validate:"required"`
	Img       []byte `validate:"required"`
}

var photordsPool *yiyidb.Kvdb

func PhotoConnect(dir string) {
	db, err := yiyidb.OpenKvdb(dir+"/cp5phts", false, false, 10) //path, enable ttl
	if err != nil {
		panic(err)
	}
	photordsPool = db
}

func delPhotos(k string) {
	vs, err := PhotostartWith(k)
	if err != nil {
		return
	}
	for _, v := range vs {
		Photodel(v)
	}
}

func photoCreate(k string, dp *PhotoDP) error {
	return photordsPool.PutJson([]byte(k), dp, 0)
}

func PhotostartWith(k string) ([]string, error) {
	ks := photordsPool.KeyStartKeys([]byte(k))
	return ks, nil
}

func PhotomaxGet(k string) (*PhotoDP, error) {
	ks := photordsPool.KeyStartKeys([]byte(k))
	if len(ks) <= 0 {
		return nil, errors.New("no data")
	}
	sortutil.Desc(ks)
	var dp PhotoDP
	err := photordsPool.GetJson([]byte(ks[0]), &dp)
	if err != nil {
		return nil, err
	}
	return &dp, nil
}

func PhotogetOneByKey(k string) (*PhotoDP, error) {
	var dp PhotoDP
	err := photordsPool.GetJson([]byte(k), &dp)
	if err != nil {
		return nil, err
	}
	return &dp, nil
}

func Photodel(k string) error {
	if len(strings.TrimSpace(k)) == 0 {
		return errors.New("key nil")
	}
	return photordsPool.Del([]byte(k))
}

func PhotoGetRange(start string, end string, interval float64, page int) ([]string, error) {
	data, err := photordsPool.KeyRange([]byte(start), []byte(end))
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
	return IntervalData, nil
}

func PhotoAll() ([]string, error) {
	ks := photordsPool.AllKeys()
	return ks, nil
}