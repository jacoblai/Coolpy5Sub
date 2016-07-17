package Nodes

import (
	"reflect"
	"github.com/garyburd/redigo/redis"
	"Coolpy/Incr"
	"encoding/json"
	"strconv"
	"Coolpy/Controller"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

type Node struct {
	Id    int64
	HubId int64 `validate:"required"`
	Title string `validate:"required"`
	About string
	Tabs  []string
	Type  int `validate:"required"`
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

var rds redis.Conn

func Connect(addr string, pwd string) {
	c, err := redis.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	_, err = c.Do("AUTH", pwd)
	if err != nil {
		panic(err)
	}
	rds = c
	rds.Do("SELECT", "3")
}

func NodeCreate(ukey string, node *Node) error {
	v, err := Incr.NodeInrc()
	if err != nil {
		return err
	}
	node.Id = v
	json, err := json.Marshal(node)
	if err != nil {
		return err
	}
	key := ukey + ":" + strconv.FormatInt(node.HubId, 10) + ":" + strconv.FormatInt(node.Id, 10)
	_, err = rds.Do("SET", key, json)
	if err != nil {
		return err
	}
	//初始化控制器
	if node.Type == NodeTypeEnum.Switcher {
		err := Controller.BeginSwitcher(ukey, node.HubId, node.Id)
		if err !=nil{
			return  errors.New("init error")
		}
	}else if node.Type == NodeTypeEnum.GenControl{
		err := Controller.BeginGenControl(ukey, node.HubId, node.Id)
		if err !=nil{
			return  errors.New("init error")
		}
	}else if node.Type == NodeTypeEnum.RangeControl{
		err := Controller.BeginRangeControl(ukey, node.HubId, node.Id)
		if err !=nil{
			return  errors.New("init error")
		}
	}
	return nil
}

//func NodeGetOne(k string) (*Node, error) {
//	o, err := redis.String(rds.Do("GET", k))
//	if err != nil {
//		return nil, err
//	}
//	h := &Node{}
//	json.Unmarshal([]byte(o), &h)
//	return h, nil
//}

func NodeStartWith(k string) ([]*Node, error) {
	data, err := redis.Strings(rds.Do("KEYSSTART", k))
	if err != nil {
		return nil, err
	}
	var ndata []*Node
	for _, v := range data {
		o, _ := redis.String(rds.Do("GET", v))
		h := &Node{}
		json.Unmarshal([]byte(o), &h)
		ndata = append(ndata, h)
	}
	return ndata, nil
}