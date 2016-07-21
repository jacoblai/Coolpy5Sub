package Account

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"encoding/json"
	"gopkg.in/go-playground/validator.v8"
	"fmt"
	"Coolpy/Deller"
)

var validate *validator.Validate

func init() {
	config := &validator.Config{TagName: "validate"}
	validate = validator.New(config)
}

func CreateAdmin() {
	if u, _ := Get("admin"); u == nil {
		p := New()
		p.Pwd = "admin"
		p.Uid = "admin"
		p.CreateUkey()
		p.UserName = "admin"
		createOrReplace(p)
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
	p.CreateUkey()
	errs := validate.Struct(p)
	if errs != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "invalid")
		return
	}
	err = createOrReplace(&p)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
		return
	}
	pStr, _ := json.Marshal(&p)
	fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
}

func UserGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	p, err := Get(ps.ByName("uid"))
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
	op, err := Get(ps.ByName("uid"))
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params nuid")
		return
	}
	if p.Uid != op.Uid {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "uidne")
		return
	}
	op.UserName = p.UserName
	op.Pwd = p.Pwd
	createOrReplace(op)
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
	del(uid)
	fmt.Fprintf(w, `{"ok":%d}`, 1)
}

func UserAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	v, err := r.Cookie("islogin")
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "dosn't login")
		return
	}
	if v.Value != "admin" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "dosn't Admin")
		return
	}
	ndata, err := all()
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
		return
	}
	pStr, _ := json.Marshal(&ndata)
	fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
}

func UserApiKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	v, err := r.Cookie("islogin")
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "dosn't login")
		return
	}
	p, err := Get(v.Value)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
		return
	}
	fmt.Fprintf(w, `{"ok":%d,"data":"%v"}`, 1, p.Ukey)
}

func UserNewApiKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	v, err := r.Cookie("islogin")
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "dosn't login")
		return
	}
	p, err := Get(v.Value)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
		return
	}
	uc, _ := r.Cookie("ukey")
	//delete all sub hub and node
	go func() {
		Deller.DelHubs <- uc.Value
	}()
	p.CreateUkey()
	createOrReplace(p)
	fmt.Fprintf(w, `{"ok":%d,"data":"%v"}`, 1, p.Ukey)
}