package Nodes

import (
	"gopkg.in/go-playground/validator.v8"
	"fmt"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"strconv"
	"Coolpy/Hubs"
	"Coolpy/Deller"
)

var validate *validator.Validate

func init() {
	config := &validator.Config{TagName: "validate"}
	validate = validator.New(config)
}

func NodePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	//post接口允许模拟put提交
	//node节点put api/hub/:hid/nodes?method=put&nid=3
	qs := r.URL.Query()
	if qs.Get("method") == "put" {
		if qs.Get("nid") != "" {
			nps := append(ps, httprouter.Param{"nid", qs.Get("nid")})
			NodePut(w, r, nps)
			return
		}
	}
	hid := ps.ByName("hid")
	if hid == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params err")
		return
	}
	decoder := json.NewDecoder(r.Body)
	var n Node
	err := decoder.Decode(&n)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
		return
	}
	_, err = r.Cookie("islogin")
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "dosn't login")
		return
	}
	nhid, err := strconv.ParseInt(hid, 10, 64)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
		return
	}
	n.HubId = nhid
	errs := validate.Struct(n)
	if errs != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "invalid")
		return
	}
	ukey, _ := r.Cookie("ukey")
	_, err = Hubs.HubGetOne(ukey.Value + ":" + hid)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext")
		return
	}
	err = nodeCreate(ukey.Value, &n)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
		return
	}
	pStr, _ := json.Marshal(&n)
	fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
}

func NodesGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	hid := ps.ByName("hid")
	if hid == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params err")
		return
	}
	_, err := r.Cookie("islogin")
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "dosn't login")
		return
	}
	ukey, _ := r.Cookie("ukey")
	ndata, err := nodeStartWith(ukey.Value + ":" + hid + ":")
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext or node not in hub")
		return
	}
	pStr, _ := json.Marshal(&ndata)
	fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
}

func NodeGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	//get接口允许模拟delete提交
	//node节点put api/hub/:hid/node/:nid?method=delete
	qs := r.URL.Query()
	if qs.Get("method") == "delete" {
		NodeDel(w, r, ps)
		return
	}
	hid := ps.ByName("hid")
	if hid == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params err")
		return
	}
	nid := ps.ByName("nid")
	if nid == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params err")
		return
	}
	_, err := r.Cookie("islogin")
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "dosn't login")
		return
	}
	ukey, _ := r.Cookie("ukey")
	ndata, err := NodeGetOne(ukey.Value + ":" + hid + ":" + nid)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext or node not in hub")
		return
	}
	pStr, _ := json.Marshal(&ndata)
	fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
}

func NodePut(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	hid := ps.ByName("hid")
	if hid == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params err")
		return
	}
	nid := ps.ByName("nid")
	if nid == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params err")
		return
	}
	decoder := json.NewDecoder(r.Body)
	var n Node
	err := decoder.Decode(&n)
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
	on, err := NodeGetOne(ukey.Value + ":" + hid + ":" + nid)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext or node not in hub")
		return
	}
	on.About = n.About
	on.Title = n.Title
	on.Tabs = n.Tabs
	nodeReplace(ukey.Value + ":" + hid + ":" + nid, on)
	pStr, _ := json.Marshal(&on)
	fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
}

func NodeDel(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	hid := ps.ByName("hid")
	if hid == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params err")
		return
	}
	nid := ps.ByName("nid")
	if nid == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params err")
		return
	}
	_, err := r.Cookie("islogin")
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "dosn't login")
		return
	}
	ukey, _ := r.Cookie("ukey")
	key := ukey.Value + ":" + hid + ":" + nid
	oh, err := NodeGetOne(key)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub nrole")
		return
	}
	if oh == nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext")
		return
	}
	//delete all sub node
	go delone(key)
	fmt.Fprintf(w, `{"ok":%d}`, 1)
}

func delone(ukeyhidnid string) {
	go func(){
		Deller.DelNode <- ukeyhidnid
	}()
	go func(){
		Deller.DelControl <- ukeyhidnid
	}()
}
