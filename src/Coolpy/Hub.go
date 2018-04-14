package Coolpy

import (
	"strconv"
	"fmt"
	"github.com/jacoblai/yiyidb"
	"strings"
	"errors"
)

type Hub struct {
	Id     uint64
	Ukey   string
	Title  string `validate:"required"`
	About  string
	Tags   []string
	Public bool
	//Local     string `validate:"required"`
	//Latitude  float64 `validate:"gte=-90,lte=90"`
	//Longitude float64 `validate:"gte=-180,lte=180"`
}

var hubrdsPool *yiyidb.Kvdb

func HubConnect(dir string) {
	db, err := yiyidb.OpenKvdb(dir+"/cp5hubs", false, false, 10) //path, enable ttl
	if err != nil {
		panic(err)
	}
	hubrdsPool = db
}

func hubCreate(hub *Hub) error {
	v, err := HubInrc()
	if err != nil {
		return err
	}
	hub.Id = v
	hub.Public = false
	key := hub.Ukey + ":" + strconv.FormatUint(hub.Id, 10)
	err = hubrdsPool.PutJson([]byte(key), &hub, 0)
	if err != nil {
		return err
	}
	return nil
}

func hubStartWith(k string) ([]*Hub, error) {
	h := Hub{}
	data, err := hubrdsPool.KeyStartByObject([]byte(k), h)
	if err != nil {
		return nil, err
	}
	var ndata []*Hub
	for _, v := range data {
		ndata = append(ndata, v.Object.(*Hub))
	}
	return ndata, nil
}

func HubGetOne(k string) (*Hub, error) {
	h := Hub{}
	err := hubrdsPool.GetJson([]byte(k), &h)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func hubReplace(k string, hub *Hub) error {
	return hubrdsPool.PutJson([]byte(k), hub, 0)
}

func delhubs(ukey string) {
	hs, err := hubStartWith(ukey)
	if err != nil {
		fmt.Println("del hub err", err)
		return
	}
	for _, v := range hs {
		k := ukey + ":" + strconv.FormatUint(v.Id, 10)
		hubdel(k)
	}
}

func hubdel(k string) error {
	if len(strings.TrimSpace(k)) == 0 {
		return errors.New("key nil")
	}

	deldos(k)

	err := hubrdsPool.Del([]byte(k))
	if err != nil {
		return err
	}
	return nil
}

func deldos(ukeyhid string) {
	DelControls(ukeyhid)
	delnodes(ukeyhid)
}

func HubAll() ([]string, error) {
	ks := hubrdsPool.AllKeys()
	return ks, nil
}
