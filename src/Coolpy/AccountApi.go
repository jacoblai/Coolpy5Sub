package Coolpy

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"encoding/json"
	"fmt"
)

func CreateAdmin() {
	if u, _ := AccGet("admin"); u == nil {
		p := AccNew()
		p.Pwd = "admin"
		p.Uid = "admin"
		p.UserName = "SuperAdmin"
		p.Email = "SuperAdmin@icoolpy.com"
		p.CreateUkey()
		Acccreate(p)
	}
}

func UserPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var p Person
	err := decoder.Decode(&p)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
		return
	}
	v, err := r.Cookie("islogin")
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "dosn't login")
		return
	}
	if v.Value != "admin" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "dosn't Admin")
		return
	}
	if p.Uid == "admin" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "admin uid")
		return
	}
	_, err = AccGet(p.Uid)
	if err == nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "uid ext")
		return
	}
	if !ValidateUidPwd(p.Uid) || !ValidateUidPwd(p.Pwd) {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "invalid Uid or Pwd")
		return
	}
	p.CreateUkey()
	err = CpValidate.Struct(p)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
		return
	}
	err = Acccreate(&p)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
		return
	}
	pStr, _ := json.Marshal(&p)
	fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
}

func UserGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	p, err := AccGet(ps.ByName("uid"))
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
		return
	}
	p.Pwd = ""
	pStr, _ := json.Marshal(&p)
	fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
}

func UserPut(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var p Person
	err := decoder.Decode(&p)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
		return
	}
	op, err := AccGet(ps.ByName("uid"))
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params nuid")
		return
	}
	if p.Uid != op.Uid {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "uidne")
		return
	}
	if !ValidateUidPwd(p.Pwd) {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "invalid")
		return
	}
	if p.UserName != "" {
		op.UserName = p.UserName
	}
	if p.Pwd != "" {
		op.Pwd = p.Pwd
	}
	if p.Email != "" {
		op.Email = p.Email
	}
	err = CpValidate.Struct(op)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
		return
	}
	Accdel(p.Uid)
	Acccreate(op)
	pStr, _ := json.Marshal(op)
	fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
}

func UserDel(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid := ps.ByName("uid")
	if uid == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params nuid")
		return
	}
	if uid == "admin" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "admin account")
		return
	}
	v, err := r.Cookie("islogin")
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "dosn't login")
		return
	}
	if v.Value != "admin" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "dosn't Admin")
		return
	}
	Accdel(uid)
	fmt.Fprintf(w, `{"ok":%d}`, 1)
}

func UserAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	v, err := r.Cookie("islogin")
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "dosn't login")
		return
	}
	if v.Value != "admin" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "dosn't Admin")
		return
	}
	ndata := AccAll()
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
		return
	}
	pStr, _ := json.Marshal(&ndata)
	fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
}

func UserApiKey(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	v, err := r.Cookie("islogin")
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "dosn't login")
		return
	}
	p, err := AccGet(v.Value)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
		return
	}
	fmt.Fprintf(w, `{"ok":%d,"data":"%v"}`, 1, p.Ukey)
}

func UserNewApiKey(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	v, err := r.Cookie("islogin")
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "dosn't login")
		return
	}
	p, err := AccGet(v.Value)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
		return
	}
	uc, _ := r.Cookie("ukey")
	//delete all sub hub and node
	delhubs(uc.Value)
	Accdel(p.Uid)
	p.CreateUkey()
	Acccreate(p)
	fmt.Fprintf(w, `{"ok":%d,"data":"%v"}`, 1, p.Ukey)
}