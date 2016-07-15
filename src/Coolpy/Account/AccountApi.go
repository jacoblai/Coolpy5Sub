package Account

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"encoding/json"
	"gopkg.in/go-playground/validator.v8"
	"fmt"
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
		CreateOrReplace(p)
	}
}

func UserPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var p Person
	err := decoder.Decode(&p)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":%v}`, 0, err)
		return
	}
	if p.Uid == "admin" {
		fmt.Fprintf(w, `{"ok":%d,"err":%v}`, 0, "isAdmin")
		return
	}
	p.CreateUkey()
	errs := validate.Struct(p)
	if errs != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":%v}`, 0, "invalid")
		return
	}
	err = CreateOrReplace(&p)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":%v}`, 0, err)
		return
	}
	pStr, _ := json.Marshal(&p)
	fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
}

func UserGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	p, err := Get(ps.ByName("uid"))
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":%v}`, 0, err)
		return
	}
	pStr, _ := json.Marshal(&p)
	fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
}

func UserPut(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var p Person
	err := decoder.Decode(&p)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":%v}`, 0, err)
		return
	}
	if p.Uid == "admin" {
		fmt.Fprintf(w, `{"ok":%d,"err":%v}`, 0, "isAdmin")
		return
	}
	op, _ := Get(ps.ByName("uid"))
	if op == nil {
		fmt.Fprintf(w, `{"ok":%d,"err":%v}`, 0, "params nuid")
		return
	}
	if p.Uid != op.Uid {
		fmt.Fprintf(w, `{"ok":%d,"err":%v}`, 0, "uidne")
		return
	}
	op.UserName = p.UserName
	op.Pwd = p.Pwd
	op.Uid = p.Uid
	CreateOrReplace(op)
	pStr, _ := json.Marshal(op)
	fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
}

func UserDel(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	p := ps.ByName("uid")
	if p == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":%v}`, 0, "params nuid")
		return
	}
	Delete(p)
	fmt.Fprintf(w, `{"ok":%d}`, 1)
}