package Coolpy

import (
	"net/http"
	"strings"
	"bytes"
	"encoding/base64"
	"github.com/julienschmidt/httprouter"
)

func Auth(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		const basicAuthPrefix string = "Basic "
		// Get the Basic Authentication credentials
		auth := r.Header.Get("Authorization")
		if strings.HasPrefix(auth, basicAuthPrefix) {
			// Check credentials
			payload, err := base64.StdEncoding.DecodeString(auth[len(basicAuthPrefix):])
			if err == nil && bytes.Contains(payload, []byte(":")) {
				pair := bytes.SplitN(payload, []byte(":"), 2)
				p, err := AccGet(string(pair[0]))
				if len(pair) == 2 && err == nil && p.Pwd == string(pair[1]) {
					r.AddCookie(&http.Cookie{
						Name:  "islogin",
						Value: p.Uid,
					})
					r.AddCookie(&http.Cookie{
						Name:  "ukey",
						Value: p.Ukey,
					})
					next(w, r, ps)
					return
				}
			}
		}
		// Request Basic Authentication otherwise
		w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}