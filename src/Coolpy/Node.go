package Coolpy

import (
	"reflect"
	"strconv"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"strings"
	"github.com/jacoblai/yiyidb"
)

type Node struct {
	Id    uint64
	HubId uint64 `validate:"required"`
	Title string `validate:"required"`
	About string
	Tags  []string
	Type  int    `validate:"required"`
	Meta  RangeMeta
}

type NodeType struct {
	Switcher     int
	GenControl   int
	RangeControl int
	Value        int
	Gps          int
	Gen          int
	Photo        int
}

var NodeTypeEnum = NodeType{1, 2, 3, 4, 5, 6, 7}
var NodeReflectType = reflect.TypeOf(NodeTypeEnum)

func (c *NodeType) GetName(v int) string {
	return NodeReflectType.Field(v).Name
}

var noderdsPool *yiyidb.Kvdb

func NodeConnect(dir string) {
	db, err := yiyidb.OpenKvdb(dir+"/cp5nodes", false, false, 10) //path, enable ttl
	if err != nil {
		panic(err)
	}
	noderdsPool = db
}

func nodeCreate(ukey string, node *Node) error {
	v, err := NodeInrc()
	if err != nil {
		return err
	}
	node.Id = v
	key := ukey + ":" + strconv.FormatUint(node.HubId, 10) + ":" + strconv.FormatUint(node.Id, 10)
	err = noderdsPool.PutJson([]byte(key), node, 0)
	if err != nil {
		return err
	}
	//验证nodetype
	if NodeTypeEnum.GetName(node.Type-1) == "" {
		return errors.New("node type error")
	}
	//初始化控制器
	if node.Type == NodeTypeEnum.Switcher {
		err := BeginSwitcher(ukey, node.HubId, node.Id)
		if err != nil {
			return errors.New("init error")
		}
	} else if node.Type == NodeTypeEnum.GenControl {
		err := BeginGenControl(ukey, node.HubId, node.Id)
		if err != nil {
			return errors.New("init error")
		}
	} else if node.Type == NodeTypeEnum.RangeControl {
		err := BeginRangeControl(ukey, node.HubId, node.Id, node.Meta)
		if err != nil {
			return errors.New("init error")
		}
	}
	return nil
}

func NodeStartWith(k string) ([]*Node, error) {
	h := Node{}
	data, err := noderdsPool.KeyStartByObject([]byte(k), h)
	if err != nil {
		return nil, err
	}
	var ndata []*Node
	for _, v := range data {
		ndata = append(ndata, v.Object.(*Node))
	}
	return ndata, nil
}

func NodeGetOne(k string) (*Node, error) {
	h := Node{}
	err := noderdsPool.GetJson([]byte(k), &h)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func nodeReplace(k string, h *Node) error {
	return noderdsPool.PutJson([]byte(k), h, 0)
}

func delnodes(ukeyhid string) {
	ns, err := NodeStartWith(ukeyhid + ":")
	if err != nil {
		return
	}
	for _, v := range ns {
		delnodes := ukeyhid + ":" + strconv.FormatUint(v.Id, 10)
		nodedel(delnodes)
	}
}

func nodedel(k string) error {
	if len(strings.TrimSpace(k)) == 0 {
		return errors.New("key nil")
	}

	deldodps(k)

	err := noderdsPool.Del([]byte(k))
	if err != nil {
		return err
	}
	return nil
}

func deldodps(ukeyhidnid string) {
	spstr := strings.Split(ukeyhidnid, ":")
	dpk := spstr[1] + "," + spstr[2]
	delValues(dpk)
	delGpss(dpk)
	delGens(dpk)
	delPhotos(dpk)
}

func NodeAll() ([]string, error) {
	ks := noderdsPool.AllKeys()
	return ks, nil
}
