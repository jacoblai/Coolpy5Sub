package Account

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"encoding/json"
	"gopkg.in/go-playground/validator.v8"
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
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if p.Uid == "admin" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	p.CreateUkey()
	errs := validate.Struct(p)
	if errs != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = CreateOrReplace(&p)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(p)
}

func UserGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	p, err := Get(ps.ByName("uid"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(p)
}

func UserPut(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var p Person
	err := decoder.Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if p.Uid == "admin" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	op, _ := Get(ps.ByName("uid"))
	if op == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if p.Uid != op.Uid {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	op.UserName = p.UserName
	op.Pwd = p.Pwd
	op.Uid = p.Uid
	CreateOrReplace(op)
	w.WriteHeader(http.StatusOK)
}

func UserDel(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	p := ps.ByName("uid")
	if p == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	Delete(p)
	w.WriteHeader(http.StatusOK)
}