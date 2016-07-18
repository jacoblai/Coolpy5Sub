package DataPoints

import (
	"fmt"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"gopkg.in/go-playground/validator.v8"
	"time"
	"Coolpy/Account"
	"strconv"
)

var validate *validator.Validate

func init() {
	config := &validator.Config{TagName: "validate"}
	validate = validator.New(config)
}

func DPPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	ukey := r.Header.Get("U-ApiKey")
	if ukey == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "ukey not post")
		return
	}
	b, err := Account.CheckKeyStart(ukey)
	if b == false {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "ukey not ext")
		return
	}
	decoder := json.NewDecoder(r.Body)
	var v ValueDP
	err = decoder.Decode(&v)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
		return
	}
	errs := validate.Struct(v)
	if errs != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "invalid")
		return
	}
	if v.TimeStamp.IsZero() {
		v.TimeStamp = time.Now()
	}
	v.HubId,_ = strconv.ParseInt(hid, 10, 64)
	v.NodeId,_ = strconv.ParseInt(nid, 10, 64)
	key := ukey + ":" + hid + ":" + nid + ":" + v.TimeStamp.Format(time.RFC3339Nano)
	err = dpCreate(key, &v)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
		return
	}
	pStr, _ := json.Marshal(&v)
	fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
}
