package Coolpy

import (
	"github.com/pmylund/sortutil"
	"time"
	"strings"
	"errors"
	"github.com/jacoblai/yiyidb"
)

type ValueDP struct {
	HubId     uint64
	NodeId    uint64
	TimeStamp time.Time
	Value     float64
}

var valdprdsPool *yiyidb.Kvdb

func ValdpConnect(dir string) {
	db, err := yiyidb.OpenKvdb(dir+"/cp5vdps", false, false, 10) //path, enable ttl
	if err != nil {
		panic(err)
	}
	valdprdsPool = db
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
	return valdprdsPool.PutJson([]byte(k), dp, 0)
}

func ValdpstartWith(k string) ([]string, error) {
	ks := valdprdsPool.KeyStartKeys([]byte(k))
	return ks, nil
}

func ValdpMaxGet(k string) (*ValueDP, error) {
	ks := valdprdsPool.KeyStartKeys([]byte(k))
	if len(ks) <= 0 {
		return nil, errors.New("no data")
	}
	sortutil.Desc(ks)
	var dp ValueDP
	err := valdprdsPool.GetJson([]byte(ks[0]), &dp)
	if err != nil {
		return nil, err
	}
	return &dp, nil
}

func ValdpGetOneByKey(k string) (*ValueDP, error) {
	var dp ValueDP
	err := valdprdsPool.GetJson([]byte(k), &dp)
	if err != nil {
		return nil, err
	}
	return &dp, nil
}

func ValdpReplace(k string, h *ValueDP) error {
	return valdprdsPool.PutJson([]byte(k), h, 0)
}

func ValdpDel(k string) error {
	if len(strings.TrimSpace(k)) == 0 {
		return errors.New("key nil")
	}
	return valdprdsPool.Del([]byte(k))
}

func ValdpGetRange(start string, end string, interval float64, page int) ([]*ValueDP, error) {
	var vdp ValueDP
	data, err := valdprdsPool.KeyRangeByObject([]byte(start), []byte(end), vdp)
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
			otime := strings.Split(IntervalData[len(IntervalData) - 1], ",")
			otm, _ := time.Parse(time.RFC3339Nano, otime[2])
			vtime := strings.Split(string(v.Key), ",")
			vtm, _ := time.Parse(time.RFC3339Nano, vtime[2])
			du := vtm.Sub(otm)
			if du.Seconds() >= interval {
				IntervalData = append(IntervalData, string(v.Key))
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
		var h ValueDP
		valdprdsPool.GetJson([]byte(v), &h)
		ndata = append(ndata, &h)
	}
	return ndata, nil
}

func ValdpAll() ([]string, error) {
	ks := valdprdsPool.AllKeys()
	return ks, nil
}