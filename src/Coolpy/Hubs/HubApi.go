package Hubs

import (
	"gopkg.in/go-playground/validator.v8"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"fmt"
	"Coolpy/Deller"
	"Coolpy/Nodes"
	"Coolpy/Controller"
)

var validate *validator.Validate

func init() {
	config := &validator.Config{TagName: "validate"}
	validate = validator.New(config)
}

func HubPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	//post接口允许模拟put提交
	//hub节点put api/hubs?method=put&hid=1
	qs := r.URL.Query()
	if qs.Get("method") == "put" {
		if qs.Get("hid") != "" {
			nps := append(ps, httprouter.Param{"hid", qs.Get("hid")})
			HubPut(w, r, nps)
			return
		}
	}
	decoder := json.NewDecoder(r.Body)
	var h Hub
	err := decoder.Decode(&h)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
		return
	}
	_, err = r.Cookie("islogin")
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "dosn't login")
		return
	}
	errs := validate.Struct(h)
	if errs != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "invalid")
		return
	}
	uc, _ := r.Cookie("ukey")
	h.Ukey = uc.Value
	err = hubCreate(&h)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
		return
	}
	pStr, _ := json.Marshal(&h)
	fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
}

func HubsGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	_, err := r.Cookie("islogin")
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "dosn't login")
		return
	}
	ukey, _ := r.Cookie("ukey")
	ndata, err := hubStartWith(ukey.Value)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
		return
	}
	pStr, _ := json.Marshal(&ndata)
	fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
}

type RHub struct {
	Id     int64
	Ukey   string
	Title  string
	About  string
	Tags   []string
	Public bool
	RNodes []*RNode
}

type RNode struct {
	Id     int64
	Title  string
	About  string
	Tags   []string
	Type   int
	Ctrler interface{}
}

func HubsAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	_, err := r.Cookie("islogin")
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "dosn't login")
		return
	}
	ukey, _ := r.Cookie("ukey")
	ndata, err := hubStartWith(ukey.Value)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
		return
	}
	var Rhub []*RHub
	for _, h := range ndata {
		nhub := &RHub{Id:h.Id, Ukey:h.Ukey, Title:h.Title, About:h.About, Tags:h.Tags, Public:h.Public}
		key := ukey + ":" + h.Id + ":"
		nodes, _ := Nodes.NodeStartWith(key)
		for _, n := range nodes {
			nnode := &RNode{Id:n.Id, Title:n.Title, About:n.About, Tags:n.Tags, Type:n.Type}
			if n.Type == Nodes.NodeTypeEnum.Switcher {
				nnode.Ctrler, _ = Controller.GetSwitcher(key + n.Id)
			} else if n.Type == Nodes.NodeTypeEnum.GenControl {
				nnode.Ctrler, _ = Controller.GetGenControl(key + n.Id)
			} else if n.Type == Nodes.NodeTypeEnum.RangeControl {
				nnode.Ctrler, _ = Controller.GetRangeControl(key + n.Id)
			}
			nhub.RNodes = append(nhub.RNodes, nnode)
		}
		Rhub = append(Rhub, nhub)
	}
	pStr, _ := json.Marshal(&Rhub)
	fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
}

func HubGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	//get接口允许模拟delete提交
	//hub节点put api/hub/:hid?method=delete
	qs := r.URL.Query()
	if qs.Get("method") == "delete" {
		HubDel(w, r, ps)
		return
	}
	hid := ps.ByName("hid")
	if hid == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params ukey")
		return
	}
	_, err := r.Cookie("islogin")
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "dosn't login")
		return
	}
	ukey, _ := r.Cookie("ukey")
	ndata, err := HubGetOne(ukey.Value + ":" + hid)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext")
		return
	}
	pStr, _ := json.Marshal(&ndata)
	fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
}

func HubPut(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	hid := ps.ByName("hid")
	if hid == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params ukey")
		return
	}
	decoder := json.NewDecoder(r.Body)
	var h Hub
	err := decoder.Decode(&h)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
		return
	}
	_, err = r.Cookie("islogin")
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "dosn't login")
		return
	}
	ukey, _ := r.Cookie("ukey")
	oh, err := HubGetOne(ukey.Value + ":" + hid)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub nrole")
		return
	}
	oh.About = h.About
	//oh.Latitude = h.Latitude
	//oh.Local = h.Local
	//oh.Longitude = h.Longitude
	oh.Public = h.Public
	oh.Tags = h.Tags
	oh.Title = h.Title
	hubReplace(oh)
	pStr, _ := json.Marshal(&oh)
	fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
}

func HubDel(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	hid := ps.ByName("hid")
	if hid == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params ukey")
		return
	}
	_, err := r.Cookie("islogin")
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "dosn't login")
		return
	}
	ukey, _ := r.Cookie("ukey")
	key := ukey.Value + ":" + hid
	_, err = HubGetOne(key)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext")
		return
	}
	//delete all sub node
	go func() {
		Deller.DelHub <- key
	}()
	fmt.Fprintf(w, `{"ok":%d}`, 1)
}