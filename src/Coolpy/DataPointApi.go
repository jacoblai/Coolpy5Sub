package Coolpy

import (
	"fmt"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
	"strconv"
	"Coolpy/Mtsvc"
)

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
			nps := append(ps, httprouter.Param{"key", qs.Get("key")})
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
	_, err := AccGetUkeyFromDb(ukey)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "ukey not ext")
		return
	}
	key := ukey + ":" + hid + ":" + nid
	dpkey := hid + "," + nid
	n, err := NodeGetOne(key)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext or node not in hub")
		return
	}
	if n.Type == NodeTypeEnum.Value {
		decoder := json.NewDecoder(r.Body)
		var v ValueDP
		err = decoder.Decode(&v)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		errs := CpValidate.Struct(v)
		if errs != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "invalid")
			return
		}
		if v.TimeStamp.IsZero() {
			v.TimeStamp = time.Now().UTC().Add(time.Hour * 8)
		}
		v.HubId, _ = strconv.ParseInt(hid, 10, 64)
		v.NodeId, _ = strconv.ParseInt(nid, 10, 64)
		err = ValueCreate(dpkey + "," + v.TimeStamp.Format(time.RFC3339Nano), &v)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&v)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == NodeTypeEnum.Gps {
		decoder := json.NewDecoder(r.Body)
		var v GpsDP
		err = decoder.Decode(&v)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		errs := CpValidate.Struct(v)
		if errs != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "invalid")
			return
		}
		if v.TimeStamp.IsZero() {
			v.TimeStamp = time.Now().UTC().Add(time.Hour * 8)
		}
		v.HubId, _ = strconv.ParseInt(hid, 10, 64)
		v.NodeId, _ = strconv.ParseInt(nid, 10, 64)
		err = GpsCreate(dpkey + "," + v.TimeStamp.Format(time.RFC3339Nano), &v)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&v)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == NodeTypeEnum.Gen {
		decoder := json.NewDecoder(r.Body)
		var v GenDP
		err = decoder.Decode(&v)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		errs := CpValidate.Struct(v)
		if errs != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "invalid")
			return
		}
		if v.TimeStamp.IsZero() {
			v.TimeStamp = time.Now().UTC().Add(time.Hour * 8)
		}
		v.HubId, _ = strconv.ParseInt(hid, 10, 64)
		v.NodeId, _ = strconv.ParseInt(nid, 10, 64)
		err = GenCreate(dpkey + "," + v.TimeStamp.Format(time.RFC3339Nano), &v)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
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
	_, err := AccGetUkeyFromDb(ukey)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "ukey not ext")
		return
	}
	key := ukey + ":" + hid + ":" + nid
	dpkey := hid + "," + nid
	n, err := NodeGetOne(key)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext or node not in hub")
		return
	}
	if n.Type == NodeTypeEnum.Value {
		max, err := ValdpMaxGet(dpkey + ",")
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&max)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == NodeTypeEnum.Gps {
		max, err := GpsdpMaxGet(dpkey + ",")
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&max)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == NodeTypeEnum.Gen {
		max, err := GendpMaxGet(dpkey + ",")
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&max)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == NodeTypeEnum.Switcher {
		c, err := GetSwitcher(key)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&c)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == NodeTypeEnum.RangeControl {
		c, err := GetRangeControl(key)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&c)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == NodeTypeEnum.GenControl {
		c, err := GetGenControl(key)
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
	_, err := AccGetUkeyFromDb(ukey)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "ukey not ext")
		return
	}
	key := ukey + ":" + hid + ":" + nid
	n, err := NodeGetOne(key)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext or node not in hub")
		return
	}
	if n.Type == NodeTypeEnum.Switcher {
		decoder := json.NewDecoder(r.Body)
		var v Switcher
		err = decoder.Decode(&v)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		errs := CpValidate.Struct(v)
		if errs != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "invalid")
			return
		}
		c, err := GetSwitcher(key)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		c.Svalue = v.Svalue
		err = ReplaceSwitcher(key, c)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&c)
		go Mtsvc.Public(key, pStr)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == NodeTypeEnum.RangeControl {
		decoder := json.NewDecoder(r.Body)
		var v RangeControl
		err = decoder.Decode(&v)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		errs := CpValidate.Struct(v)
		if errs != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "invalid")
			return
		}
		c, err := GetRangeControl(key)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		if v.Rvalue > c.Max || c.Rvalue < c.Min {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "range value err")
			return
		}
		c.Rvalue = v.Rvalue
		err = ReplaceRangeControl(key, c)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&c)
		go Mtsvc.Public(key, pStr)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == NodeTypeEnum.GenControl {
		decoder := json.NewDecoder(r.Body)
		var v GenControl
		err = decoder.Decode(&v)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		errs := CpValidate.Struct(v)
		if errs != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "invalid")
			return
		}
		c, err := GetGenControl(key)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		c.Gvalue = v.Gvalue
		err = ReplaceGenControl(key, c)
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
	//get接口允许模拟delete提交
	//模拟控制器put api/hub/:hid/node/:nid/datapoint/:key?method=delete
	qs := r.URL.Query()
	if qs.Get("method") == "delete" {
		DPDelByKey(w, r, ps)
		return
	}
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
	_, err := AccGetUkeyFromDb(ukey)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "ukey not ext")
		return
	}
	key := ukey + ":" + hid + ":" + nid
	dpkey := hid + "," + nid
	n, err := NodeGetOne(key)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext or node not in hub")
		return
	}
	if n.Type == NodeTypeEnum.Value {
		one, err := ValdpGetOneByKey(dpkey + "," + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&one)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == NodeTypeEnum.Gps {
		one, err := GpsdpGetOneByKey(dpkey + "," + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&one)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == NodeTypeEnum.Gen {
		one, err := GendpGetOneByKey(dpkey + "," + dpKey)
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
	_, err := AccGetUkeyFromDb(ukey)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "ukey not ext")
		return
	}
	key := ukey + ":" + hid + ":" + nid
	dpkey := hid + "," + nid
	n, err := NodeGetOne(key)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext or node not in hub")
		return
	}
	if n.Type == NodeTypeEnum.Value {
		decoder := json.NewDecoder(r.Body)
		var v ValueDP
		err = decoder.Decode(&v)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		errs := CpValidate.Struct(v)
		if errs != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "invalid")
			return
		}
		c, err := ValdpGetOneByKey(dpkey + "," + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		c.Value = v.Value
		err = ValdpReplace(dpkey + "," + dpKey, c)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&c)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == NodeTypeEnum.Gps {
		decoder := json.NewDecoder(r.Body)
		var v GpsDP
		err = decoder.Decode(&v)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		errs := CpValidate.Struct(v)
		if errs != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "invalid")
			return
		}
		c, err := GpsdpGetOneByKey(dpkey + "," + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		c.Lat = v.Lat
		c.Lng = v.Lng
		c.Offset = v.Offset
		c.Speed = v.Speed
		err = GpsdpReplace(dpkey + "," + dpKey, c)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&c)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == NodeTypeEnum.Gen {
		decoder := json.NewDecoder(r.Body)
		var v GenDP
		err = decoder.Decode(&v)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		errs := CpValidate.Struct(v)
		if errs != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "invalid")
			return
		}
		c, err := GendpGetOneByKey(dpkey + "," + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		c.Value = v.Value
		err = GendpReplace(dpkey + "," + dpKey, c)
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
	_, err := AccGetUkeyFromDb(ukey)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "ukey not ext")
		return
	}
	key := ukey + ":" + hid + ":" + nid
	dpkey := hid + "," + nid
	n, err := NodeGetOne(key)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext or node not in hub")
		return
	}
	if n.Type == NodeTypeEnum.Value {
		c, err := ValdpGetOneByKey(dpkey + "," + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		err = ValdpDel(dpkey + "," + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&c)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == NodeTypeEnum.Gps {
		c, err := GpsdpGetOneByKey(dpkey + "," + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		err = GpsdpDel(dpkey + "," + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&c)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == NodeTypeEnum.Gen {
		c, err := GendpGetOneByKey(dpkey + "," + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		err = GendpDel(dpkey + "," + dpKey)
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
	_, err = AccGetUkeyFromDb(ukey)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "ukey not ext")
		return
	}
	key := ukey + ":" + hid + ":" + nid
	dpkey := hid + "," + nid
	n, err := NodeGetOne(key)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext or node not in hub")
		return
	}
	if n.Type == NodeTypeEnum.Value {
		c, err := ValdpGetRange(dpkey + "," + start.Format(time.RFC3339Nano),
			dpkey + "," + end.Format(time.RFC3339Nano), interval, page)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&c)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == NodeTypeEnum.Gps {
		c, err := GpsdpGetRange(dpkey + "," + start.Format(time.RFC3339Nano),
			dpkey + "," + end.Format(time.RFC3339Nano), interval, page)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&c)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == NodeTypeEnum.Gen {
		c, err := GendpGetRange(dpkey + "," + start.Format(time.RFC3339Nano),
			dpkey + "," + end.Format(time.RFC3339Nano), interval, page)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&c)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else if n.Type == NodeTypeEnum.Photo {
		c, err := PhotoGetRange(dpkey + "," + start.Format(time.RFC3339Nano),
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