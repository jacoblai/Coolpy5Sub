package Photos

import (
	"fmt"
	"strconv"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"Coolpy/Account"
	"Coolpy/Nodes"
	"time"
	"io/ioutil"
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
	b, err := Account.CheckUKey(ukey + ":")
	if b == false {
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
	dpkey := ukey + "," + hid + "," + nid
	n, err := Nodes.NodeGetOne(key)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "hub not ext or node not in hub")
		return
	}
	if n.Type == Nodes.NodeTypeEnum.Photo {
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
		p.HubId, _ = strconv.ParseInt(hid, 10, 64)
		p.NodeId, _ = strconv.ParseInt(nid, 10, 64)
		p.Img = img
		p.Mime = mm
		err = PhotoCreate(dpkey + "," + p.TimeStamp.Format(time.RFC3339Nano), &p)
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		pStr, _ := json.Marshal(&p)
		fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, string(pStr))
	} else {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "unkown type")
	}
}
