package Coolpy

import (
	"strconv"
	"strings"
	"errors"
	"github.com/jacoblai/yiyidb"
)

type Switcher struct {
	HubId  uint64
	NodeId uint64
	Svalue int
}

type GenControl struct {
	HubId  uint64
	NodeId uint64
	Gvalue string `validate:"required"`
}

type RangeControl struct {
	HubId  uint64
	NodeId uint64
	Rvalue int64
	Min    int64
	Max    int64
	Step   int64
}

type RangeMeta struct {
	Min  int64
	Max  int64
	Step int64
}

var ctrlrdsPool *yiyidb.Kvdb

func CtrlConnect(dir string) {
	db, err := yiyidb.OpenKvdb(dir+"/cp5ctrls", false, false, 10) //path, enable ttl
	if err != nil {
		panic(err)
	}
	ctrlrdsPool = db
}

func DelControls(k string) {
	cs, err := ctrlStartWith(k)
	if err != nil {
		return
	}
	for _, v := range cs {
		ctrldel(v)
	}
}

func ReplaceSwitcher(k string, s *Switcher) error {
	return ctrlrdsPool.PutJson([]byte(k), s, 0)
}

func GetSwitcher(k string) (*Switcher, error) {
	var s Switcher
	err := ctrlrdsPool.GetJson([]byte(k), &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func BeginSwitcher(ukey string, Hubid uint64, Nodeid uint64) error {
	key := ukey + ":" + strconv.FormatUint(Hubid, 10) + ":" + strconv.FormatUint(Nodeid, 10)
	o := Switcher{
		HubId:  Hubid,
		NodeId: Nodeid,
		Svalue: 0,
	}
	return ctrlrdsPool.PutJson([]byte(key), &o, 0)
}

func ReplaceRangeControl(k string, s *RangeControl) error {
	return ctrlrdsPool.PutJson([]byte(k), s, 0)
}

func GetRangeControl(k string) (*RangeControl, error) {
	var r RangeControl
	err := ctrlrdsPool.GetJson([]byte(k), &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func BeginRangeControl(ukey string, Hubid uint64, Nodeid uint64, meta RangeMeta) error {
	key := ukey + ":" + strconv.FormatUint(Hubid, 10) + ":" + strconv.FormatUint(Nodeid, 10)
	if meta.Max == 0 {
		meta.Max = 255
	}
	if meta.Step == 0 {
		meta.Step = 5
	}
	o := RangeControl{
		HubId:  Hubid,
		NodeId: Nodeid,
		Rvalue: 0,
		Min:    meta.Min,
		Max:    meta.Max,
		Step:   meta.Step,
	}
	return ctrlrdsPool.PutJson([]byte(key), &o, 0)
}

func ReplaceGenControl(k string, s *GenControl) error {
	return ctrlrdsPool.PutJson([]byte(k), s, 0)
}

func GetGenControl(k string) (*GenControl, error) {
	var r GenControl
	err := ctrlrdsPool.GetJson([]byte(k), &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func BeginGenControl(ukey string, Hubid uint64, Nodeid uint64) error {
	key := ukey + ":" + strconv.FormatUint(Hubid, 10) + ":" + strconv.FormatUint(Nodeid, 10)
	o := GenControl{
		HubId:  Hubid,
		NodeId: Nodeid,
		Gvalue: "",
	}
	return ctrlrdsPool.PutJson([]byte(key), &o, 0)
}

func ctrlStartWith(k string) ([]string, error) {
	ks := ctrlrdsPool.KeyStartKeys([]byte(k))
	return ks, nil
}

func ctrldel(k string) error {
	if len(strings.TrimSpace(k)) == 0 {
		return errors.New("key nil")
	}
	return ctrlrdsPool.Del([]byte(k))
}

func CtrlAll() ([]string, error) {
	ks := ctrlrdsPool.AllKeys()
	return ks, nil
}
