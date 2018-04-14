package Coolpy

import (
	"fmt"
	"strconv"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
	"io/ioutil"
	"bytes"
)

func PhotoPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	//限制上传图片大小
	l, _ := strconv.Atoi(r.Header.Get("Content-Length"))
	if l > 500 * 1024 {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "just upload less 500k")
		return
	}
	key := ukey + ":" + hid + ":" + nid
	dpkey := hid + "," + nid
	n, err := NodeGetOne(key)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext or node not in hub")
		return
	}
	if n.Type == NodeTypeEnum.Photo {
		img, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "read err")
			return
		}
		mm := mimeComput(img)
		if mm == "" {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "just png jpg gif file")
			return
		}
		var p PhotoDP
		if p.TimeStamp.IsZero() {
			p.TimeStamp = time.Now().UTC().Add(time.Hour * 8)
		}
		p.HubId, _ = strconv.ParseUint(hid, 10, 64)
		p.NodeId, _ = strconv.ParseUint(nid, 10, 64)
		p.Img = img
		p.Mime = mm
		p.Size = int64(len(p.Img))
		err = photoCreate(dpkey + "," + p.TimeStamp.Format(time.RFC3339Nano), &p)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		fmt.Fprintf(w, `{"ok":%d,"data":"%v"}`, 1, p.TimeStamp.Format(time.RFC3339Nano))
	} else {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "unkown type")
	}
}

func PhotoGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	if n.Type == NodeTypeEnum.Photo {
		max, err := PhotomaxGet(dpkey + ",")
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		w.Header().Set("Content-Type", max.Mime)
		if r.Header.Get("Range") != "" {
			f := bytes.NewReader(max.Img)
			start_byte := parseRange(r.Header.Get("Range"))
			if start_byte < max.Size {
				f.Seek(start_byte, 0)
			} else {
				start_byte = 0
			}
			w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start_byte, max.Size - 1, max.Size))
			f.WriteTo(w)
		} else {
			w.Write(max.Img)
		}
	} else {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "unkown type")
	}
}

func PhotoGetByKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
		qs := r.URL.Query()
		if qs.Get("ukey") != "" {
			ukey = qs.Get("ukey")
		} else {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "ukey not post")
			return
		}
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
	if n.Type == NodeTypeEnum.Photo {
		one, err := PhotogetOneByKey(dpkey + "," + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		w.Header().Set("Content-Type", one.Mime)
		if r.Header.Get("Range") != "" {
			f := bytes.NewReader(one.Img)
			start_byte := parseRange(r.Header.Get("Range"))
			if start_byte < one.Size {
				f.Seek(start_byte, 0)
			} else {
				start_byte = 0
			}
			w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start_byte, one.Size - 1, one.Size))
			f.WriteTo(w)
		} else {
			w.Write(one.Img)
		}
	} else {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "unkown type")
	}
}

func PhotoDelByKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	if n.Type == NodeTypeEnum.Photo {
		c, err := PhotogetOneByKey(dpkey + "," + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		err = Photodel(dpkey + "," + dpKey)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		fmt.Fprintf(w, `{"ok":%d,"data":"%v"}`, 1, c.TimeStamp.Format(time.RFC3339Nano))
	} else {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "unkown type")
	}
}
