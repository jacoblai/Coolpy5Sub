package Cors

import "net/http"

const (
	options string = "OPTIONS"
	allow_origin string = "Access-Control-Allow-Origin"
	allow_methods string = "Access-Control-Allow-Methods"
	allow_headers string = "Access-Control-Allow-Headers"
	origin string = "Origin"
	methods string = "GET,PUT,POST,DELETE"
	// If you want to expose some other headers add it here
	headers string = "U-ApiKey,Authorization,Content-Length,Content-Type"
)

// Handler will allow cross-origin HTTP requests
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set allow origin to match origin of our request or fall back to *
		if o := r.Header.Get(origin); o != "" {
			w.Header().Set(allow_origin, o)
		} else {
			w.Header().Set(allow_origin, "*")
		}

		// Set other headers
		w.Header().Set(allow_headers, headers)
		w.Header().Set(allow_methods, methods)

		// If this was preflight options request let's write empty ok response and return
		if r.Method == options {
			w.WriteHeader(http.StatusOK)
			w.Write(nil)
			return
		}
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		next.ServeHTTP(w, r)
	})
}