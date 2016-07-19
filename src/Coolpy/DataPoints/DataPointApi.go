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
	"Coolpy/Hubs"
	"Coolpy/Nodes"
	"Coolpy/Values"
	"Coolpy/Gpss"
	"Coolpy/Gens"
	"Coolpy/Controller"
	"Coolpy/Mtsvc"
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
	if k, _ := Hubs.CheckHubId(hid); k == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext")
		return
	}
	nid := ps.ByName("nid")
	if nid == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params err")
		return
	}
	if k, _ := Nodes.CheckNodeId(nid); k == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "node not ext")
		return
	}
	ukey := r.Header.Get("U-ApiKey")
	if ukey == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "ukey not post")
		return
	}
	b, err := Account.CheckUKey(ukey + ":")
	if b == false {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "ukey not ext")
		return
	}
	key := ukey + ":" + hid + ":" + nid
	n, err := Nodes.NodeGetOne(key)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext or node not in hub")
		return
	}
	if n.Type == Nodes.NodeTypeEnum.Value {
		decoder := json.NewDecoder(r.Body)
		var v Values.ValueDP
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
			v.TimeStamp = time.Now().UTC().Add(time.Hour * 8)
		}
		v.HubId, _ = strconv.ParseInt(hid, 10, 64)
		v.NodeId, _ = strconv.ParseInt(nid, 10, 64)
		err = Values.ValueCreate(key + ":" + v.TimeStamp.Format(time.RFC3339Nano), &v)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&v)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == Nodes.NodeTypeEnum.Gps {
		decoder := json.NewDecoder(r.Body)
		var v Gpss.GpsDP
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
			v.TimeStamp = time.Now().UTC().Add(time.Hour * 8)
		}
		v.HubId, _ = strconv.ParseInt(hid, 10, 64)
		v.NodeId, _ = strconv.ParseInt(nid, 10, 64)
		err = Gpss.GpsCreate(key + ":" + v.TimeStamp.Format(time.RFC3339Nano), &v)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&v)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == Nodes.NodeTypeEnum.Gen {
		decoder := json.NewDecoder(r.Body)
		var v Gens.GenDP
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
			v.TimeStamp = time.Now().UTC().Add(time.Hour * 8)
		}
		v.HubId, _ = strconv.ParseInt(hid, 10, 64)
		v.NodeId, _ = strconv.ParseInt(nid, 10, 64)
		err = Gens.GenCreate(key + ":" + v.TimeStamp.Format(time.RFC3339Nano), &v)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&v)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "unkown type")
	}
}

func DPGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	hid := ps.ByName("hid")
	if hid == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params err")
		return
	}
	if k, _ := Hubs.CheckHubId(hid); k == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext")
		return
	}
	nid := ps.ByName("nid")
	if nid == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params err")
		return
	}
	if k, _ := Nodes.CheckNodeId(nid); k == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "node not ext")
		return
	}
	ukey := r.Header.Get("U-ApiKey")
	if ukey == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "ukey not post")
		return
	}
	b, err := Account.CheckUKey(ukey + ":")
	if b == false {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "ukey not ext")
		return
	}
	key := ukey + ":" + hid + ":" + nid
	n, err := Nodes.NodeGetOne(key)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext or node not in hub")
		return
	}
	if n.Type == Nodes.NodeTypeEnum.Value {
		max, err := Values.MaxGet(key + ":")
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&max)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == Nodes.NodeTypeEnum.Gps {
		max, err := Gpss.MaxGet(key + ":")
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&max)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == Nodes.NodeTypeEnum.Gen {
		max, err := Gens.MaxGet(key + ":")
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&max)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == Nodes.NodeTypeEnum.Switcher {
		c, err := Controller.GetSwitcher(key)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&c)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == Nodes.NodeTypeEnum.RangeControl {
		c, err := Controller.GetRangeControl(key)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&c)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == Nodes.NodeTypeEnum.GenControl {
		c, err := Controller.GetGenControl(key)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&c)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "unkown type")
	}
}

func DPPut(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	hid := ps.ByName("hid")
	if hid == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params err")
		return
	}
	if k, _ := Hubs.CheckHubId(hid); k == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext")
		return
	}
	nid := ps.ByName("nid")
	if nid == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params err")
		return
	}
	if k, _ := Nodes.CheckNodeId(nid); k == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "node not ext")
		return
	}
	ukey := r.Header.Get("U-ApiKey")
	if ukey == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "ukey not post")
		return
	}
	b, err := Account.CheckUKey(ukey + ":")
	if b == false {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "ukey not ext")
		return
	}
	key := ukey + ":" + hid + ":" + nid
	n, err := Nodes.NodeGetOne(key)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext or node not in hub")
		return
	}
	if n.Type == Nodes.NodeTypeEnum.Switcher {
		decoder := json.NewDecoder(r.Body)
		var v Controller.Switcher
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
		c, err := Controller.GetSwitcher(key)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		c.Svalue = v.Svalue
		err = Controller.ReplaceSwitcher(key, c)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&c)
		Mtsvc.Public(key, pStr)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == Nodes.NodeTypeEnum.RangeControl {
		decoder := json.NewDecoder(r.Body)
		var v Controller.RangeControl
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
		c, err := Controller.GetRangeControl(key)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		if v.Rvalue > c.Max || c.Rvalue < c.Min {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "range value err")
			return
		}
		c.Rvalue = v.Rvalue
		err = Controller.ReplaceRangeControl(key, c)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&c)
		Mtsvc.Public(key, pStr)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == Nodes.NodeTypeEnum.GenControl {
		decoder := json.NewDecoder(r.Body)
		var v Controller.GenControl
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
		c, err := Controller.GetGenControl(key)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		c.Gvalue = v.Gvalue
		err = Controller.ReplaceGenControl(key, c)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&c)
		Mtsvc.Public(key, pStr)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "unkown type")
	}
}

func DPGetByKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	dpKey := ps.ByName("key")
	if dpKey == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params err")
		return
	}
	hid := ps.ByName("hid")
	if hid == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params err")
		return
	}
	if k, _ := Hubs.CheckHubId(hid); k == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext")
		return
	}
	nid := ps.ByName("nid")
	if nid == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params err")
		return
	}
	if k, _ := Nodes.CheckNodeId(nid); k == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "node not ext")
		return
	}
	ukey := r.Header.Get("U-ApiKey")
	if ukey == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "ukey not post")
		return
	}
	b, err := Account.CheckUKey(ukey + ":")
	if b == false {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "ukey not ext")
		return
	}
	key := ukey + ":" + hid + ":" + nid
	n, err := Nodes.NodeGetOne(key)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext or node not in hub")
		return
	}
	if n.Type == Nodes.NodeTypeEnum.Value {
		one, err := Values.GetOneByKey(key + ":" + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&one)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == Nodes.NodeTypeEnum.Gps {
		one, err := Gpss.GetOneByKey(key + ":" + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&one)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == Nodes.NodeTypeEnum.Gen {
		max, err := Gens.GetOneByKey(key + ":" + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&max)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "unkown type")
	}
}

func DPPutByKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	dpKey := ps.ByName("key")
	if dpKey == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params err")
		return
	}
	hid := ps.ByName("hid")
	if hid == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params err")
		return
	}
	if k, _ := Hubs.CheckHubId(hid); k == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext")
		return
	}
	nid := ps.ByName("nid")
	if nid == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params err")
		return
	}
	if k, _ := Nodes.CheckNodeId(nid); k == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "node not ext")
		return
	}
	ukey := r.Header.Get("U-ApiKey")
	if ukey == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "ukey not post")
		return
	}
	b, err := Account.CheckUKey(ukey + ":")
	if b == false {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "ukey not ext")
		return
	}
	key := ukey + ":" + hid + ":" + nid
	n, err := Nodes.NodeGetOne(key)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext or node not in hub")
		return
	}
	if n.Type == Nodes.NodeTypeEnum.Value {
		decoder := json.NewDecoder(r.Body)
		var v Values.ValueDP
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
		c, err := Values.GetOneByKey(key + ":" + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		c.Value = v.Value
		err = Values.Replace(key + ":" + dpKey, c)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&c)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == Nodes.NodeTypeEnum.Gps {
		decoder := json.NewDecoder(r.Body)
		var v Gpss.GpsDP
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
		c, err := Gpss.GetOneByKey(key + ":" + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		c.Lat = v.Lat
		c.Lng = v.Lng
		c.Offset = v.Offset
		c.Speed = v.Speed
		err = Gpss.Replace(key + ":" + dpKey, c)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&c)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == Nodes.NodeTypeEnum.Gen {
		decoder := json.NewDecoder(r.Body)
		var v Gens.GenDP
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
		c, err := Gens.GetOneByKey(key + ":" + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		c.Value = v.Value
		err = Gens.Replace(key + ":" + dpKey, c)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&c)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "unkown type")
	}
}

func DPDelByKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	dpKey := ps.ByName("key")
	if dpKey == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params err")
		return
	}
	hid := ps.ByName("hid")
	if hid == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params err")
		return
	}
	if k, _ := Hubs.CheckHubId(hid); k == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext")
		return
	}
	nid := ps.ByName("nid")
	if nid == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params err")
		return
	}
	if k, _ := Nodes.CheckNodeId(nid); k == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "node not ext")
		return
	}
	ukey := r.Header.Get("U-ApiKey")
	if ukey == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "ukey not post")
		return
	}
	b, err := Account.CheckUKey(ukey + ":")
	if b == false {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "ukey not ext")
		return
	}
	key := ukey + ":" + hid + ":" + nid
	n, err := Nodes.NodeGetOne(key)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext or node not in hub")
		return
	}
	if n.Type == Nodes.NodeTypeEnum.Value {
		c, err := Values.GetOneByKey(key + ":" + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		err = Values.Delete(key + ":" + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&c)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == Nodes.NodeTypeEnum.Gps {
		c, err := Gpss.GetOneByKey(key + ":" + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		err = Gpss.Delete(key + ":" + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&c)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == Nodes.NodeTypeEnum.Gen {
		c, err := Gens.GetOneByKey(key + ":" + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		err = Gens.Delete(key + ":" + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&c)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "unkown type")
	}
}

func DPGetRange(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	qs := r.URL.Query()
	dpStart := qs.Get("start")
	start, err := time.Parse(time.RFC3339Nano, dpStart)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params err")
		return
	}
	dpEnd := qs.Get("end")
	end, err := time.Parse(time.RFC3339Nano, dpEnd)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params err")
		return
	}
	dpInterval := qs.Get("interval")
	interval, err := strconv.Atoi(dpInterval)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params err")
		return
	}
	dpPage := qs.Get("page")
	page, err := strconv.Atoi(dpPage)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params err")
		return
	}
	hid := ps.ByName("hid")
	if hid == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params err")
		return
	}
	if k, _ := Hubs.CheckHubId(hid); k == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext")
		return
	}
	nid := ps.ByName("nid")
	if nid == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params err")
		return
	}
	if k, _ := Nodes.CheckNodeId(nid); k == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "node not ext")
		return
	}
	ukey := r.Header.Get("U-ApiKey")
	if ukey == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "ukey not post")
		return
	}
	b, err := Account.CheckUKey(ukey + ":")
	if b == false {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "ukey not ext")
		return
	}
	key := ukey + ":" + hid + ":" + nid
	n, err := Nodes.NodeGetOne(key)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext or node not in hub")
		return
	}
	if n.Type == Nodes.NodeTypeEnum.Value {
		c, err := Values.GetRange(key + ":" + start.Format(time.RFC3339Nano),
			key + ":" + end.Format(time.RFC3339Nano), interval, page)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&c)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == Nodes.NodeTypeEnum.Gps {

	} else if n.Type == Nodes.NodeTypeEnum.Gen {

	} else {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "unkown type")
	}
}