package Hubs

import (
	"gopkg.in/go-playground/validator.v8"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"fmt"
)

var validate *validator.Validate

func init() {
	config := &validator.Config{TagName: "validate"}
	validate = validator.New(config)
}

func HubPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var h Hub
	err := decoder.Decode(&h)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":%v}`, 0, err)
		return
	}
	_, err = r.Cookie("islogin")
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":%v}`, 0, "dosn't login")
		return
	}
	errs := validate.Struct(h)
	if errs != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":%v}`, 0, "invalid")
		return
	}
	uc, _ := r.Cookie("ukey")
	h.Ukey = uc.Value
	err = HubCreate(&h)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":%v}`, 0, err)
		return
	}
	pStr, _ := json.Marshal(&h)
	fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
}

func HubsGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uk := ps.ByName("ukey")
	if uk == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":%v}`, 0, "params ukey")
		return
	}
	_, err := r.Cookie("islogin")
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":%v}`, 0, "dosn't login")
		return
	}
	ukey, _ := r.Cookie("ukey")
	if ukey.Value != uk{
		fmt.Fprintf(w, `{"ok":%d,"err":%v}`, 0, "ukey inrole")
		return
	}
	ndata, err := HubStartWith(uk)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":%v}`, 0, err)
		return
	}
	pStr, _ := json.Marshal(&ndata)
	fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
}