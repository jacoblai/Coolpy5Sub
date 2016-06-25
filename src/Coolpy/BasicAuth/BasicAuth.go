package Basicauth

import (
	"net/http"
	"strings"
	"bytes"
	"encoding/base64"
	"github.com/julienschmidt/httprouter"
)

func Auth(next httprouter.Handle)httprouter.Handle  {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user := []byte("jac")
		pass := []byte("jac")

		const basicAuthPrefix string = "Basic "
		// Get the Basic Authentication credentials
		auth := r.Header.Get("Authorization")
		if strings.HasPrefix(auth, basicAuthPrefix) {
			// Check credentials
			payload, err := base64.StdEncoding.DecodeString(auth[len(basicAuthPrefix):])
			if err == nil {
				pair := bytes.SplitN(payload, []byte(":"), 2)
				if len(pair) == 2 &&
				bytes.Equal(pair[0], user) &&
				bytes.Equal(pair[1], pass) {
					// Delegate request to the given handle
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