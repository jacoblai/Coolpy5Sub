package Coolpy

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"fmt"
	"os/exec"
	"strings"
	"encoding/json"
	"strconv"
	"io/ioutil"
)

type CmdDP struct {
	Cmd string
}

func CmdPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var c CmdDP
	err := decoder.Decode(&c)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
		return
	}
	err = CpValidate.Struct(c)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
		return
	}
	parts := strings.Fields(c.Cmd)
	head := parts[0]
	parts = parts[1:]

	out, err := exec.Command(head, parts...).Output()
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":%v}`, 0, strconv.Quote(err.Error()))
	}
	fmt.Fprintf(w, `{"ok":%d,"data":%v}`, 1, strconv.Quote(string(out)))
}

func UploadPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	fn := ps.ByName("filename")
	if fn == "" {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "params err")
		return
	}
	audio, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "read err")
		return
	}
	if !IsMp3(audio[:3]) {
		fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, "must mp3")
		return
	}
	fp := "/tmp/" + fn
	ioutil.WriteFile(fp, audio, 0644)
	fmt.Fprintf(w, `{"ok":%d,"data":"%v"}`, 1, fp)
}

func IsMp3(buf []byte) bool {
	return len(buf) > 2 && ((buf[0] == 0x49 && buf[1] == 0x44 && buf[2] == 0x33) || (buf[0] == 0xFF && buf[1] == 0xfb))
}
