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
	"Coolpy/Nodes"
	"Coolpy/Values"
	"Coolpy/Gpss"
	"Coolpy/Gens"
	"Coolpy/Controller"
	"Coolpy/Mtsvc"
	"Coolpy/Photos"
)

var validate *validator.Validate

func init() {
	config := &validator.Config{TagName: "validate"}
	validate = validator.New(config)
}

func DPPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	//post接口允许模拟put提交
	//模拟控制器put api/hub/:hid/node/:nid/datapoints?method=put
	//模拟传感器put api/hub/:hid/node/:nid/datapoints?method=put&key=2015-11-12T02:10:55.5245871Z
	qs := r.URL.Query()
	if qs.Get("method") == "put" {
		if qs.Get("key") == "" {
			DPPut(w, r, ps)
			return
		} else {
			nps := httprouter.Params{
				httprouter.Param{"key", qs.Get("key")},
			}
			DPPutByKey(w, r, nps)
			return
		}
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
	ukey := r.Header.Get("U-ApiKey")
	if ukey == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "ukey not post")
		return
	}
	_, err := Account.GetUkeyFromDb(ukey)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "ukey not ext")
		return
	}
	key := ukey + ":" + hid + ":" + nid
	dpkey := ukey + "," + hid + "," + nid
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
		err = Values.ValueCreate(dpkey + "," + v.TimeStamp.Format(time.RFC3339Nano), &v)
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
		err = Gpss.GpsCreate(dpkey + "," + v.TimeStamp.Format(time.RFC3339Nano), &v)
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
		err = Gens.GenCreate(dpkey + "," + v.TimeStamp.Format(time.RFC3339Nano), &v)
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
	_, err := Account.GetUkeyFromDb(ukey)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "ukey not ext")
		return
	}
	key := ukey + ":" + hid + ":" + nid
	dpkey := ukey + "," + hid + "," + nid
	n, err := Nodes.NodeGetOne(key)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext or node not in hub")
		return
	}
	if n.Type == Nodes.NodeTypeEnum.Value {
		max, err := Values.MaxGet(dpkey + ",")
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&max)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == Nodes.NodeTypeEnum.Gps {
		max, err := Gpss.MaxGet(dpkey + ",")
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&max)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == Nodes.NodeTypeEnum.Gen {
		max, err := Gens.MaxGet(dpkey + ",")
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
	_, err := Account.GetUkeyFromDb(ukey)
	if err != nil {
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
		go Mtsvc.Public(key, pStr)
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
		go Mtsvc.Public(key, pStr)
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
		go Mtsvc.Public(key, pStr)
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
	_, err := Account.GetUkeyFromDb(ukey)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "ukey not ext")
		return
	}
	key := ukey + ":" + hid + ":" + nid
	dpkey := ukey + "," + hid + "," + nid
	n, err := Nodes.NodeGetOne(key)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext or node not in hub")
		return
	}
	if n.Type == Nodes.NodeTypeEnum.Value {
		one, err := Values.GetOneByKey(dpkey + "," + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&one)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == Nodes.NodeTypeEnum.Gps {
		one, err := Gpss.GetOneByKey(dpkey + "," + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&one)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == Nodes.NodeTypeEnum.Gen {
		one, err := Gens.GetOneByKey(dpkey + "," + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&one)
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
	_, err := Account.GetUkeyFromDb(ukey)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "ukey not ext")
		return
	}
	key := ukey + ":" + hid + ":" + nid
	dpkey := ukey + "," + hid + "," + nid
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
		c, err := Values.GetOneByKey(dpkey + "," + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		c.Value = v.Value
		err = Values.Replace(dpkey + "," + dpKey, c)
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
		c, err := Gpss.GetOneByKey(dpkey + "," + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		c.Lat = v.Lat
		c.Lng = v.Lng
		c.Offset = v.Offset
		c.Speed = v.Speed
		err = Gpss.Replace(dpkey + "," + dpKey, c)
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
		c, err := Gens.GetOneByKey(dpkey + "," + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		c.Value = v.Value
		err = Gens.Replace(dpkey + "," + dpKey, c)
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
	_, err := Account.GetUkeyFromDb(ukey)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "ukey not ext")
		return
	}
	key := ukey + ":" + hid + ":" + nid
	dpkey := ukey + "," + hid + "," + nid
	n, err := Nodes.NodeGetOne(key)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext or node not in hub")
		return
	}
	if n.Type == Nodes.NodeTypeEnum.Value {
		c, err := Values.GetOneByKey(dpkey + "," + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		err = Values.Del(dpkey + "," + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&c)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == Nodes.NodeTypeEnum.Gps {
		c, err := Gpss.GetOneByKey(dpkey + "," + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		err = Gpss.Del(dpkey + "," + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&c)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == Nodes.NodeTypeEnum.Gen {
		c, err := Gens.GetOneByKey(dpkey + "," + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		err = Gens.Del(dpkey + "," + dpKey)
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
	interval, err := strconv.ParseFloat(dpInterval, 10)
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
	_, err = Account.GetUkeyFromDb(ukey)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "ukey not ext")
		return
	}
	key := ukey + ":" + hid + ":" + nid
	dpkey := ukey + "," + hid + "," + nid
	n, err := Nodes.NodeGetOne(key)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext or node not in hub")
		return
	}
	if n.Type == Nodes.NodeTypeEnum.Value {
		c, err := Values.GetRange(dpkey + "," + start.Format(time.RFC3339Nano),
			dpkey + "," + end.Format(time.RFC3339Nano), interval, page)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&c)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == Nodes.NodeTypeEnum.Gps {
		c, err := Gpss.GetRange(dpkey + "," + start.Format(time.RFC3339Nano),
			dpkey + "," + end.Format(time.RFC3339Nano), interval, page)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&c)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == Nodes.NodeTypeEnum.Gen {
		c, err := Gens.GetRange(dpkey + "," + start.Format(time.RFC3339Nano),
			dpkey + "," + end.Format(time.RFC3339Nano), interval, page)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&c)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == Nodes.NodeTypeEnum.Photo {
		c, err := Photos.GetRange(dpkey + "," + start.Format(time.RFC3339Nano),
			dpkey + "," + end.Format(time.RFC3339Nano), interval, page)
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