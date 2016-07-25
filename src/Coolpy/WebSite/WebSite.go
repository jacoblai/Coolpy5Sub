package WebSite

import (
	"net/http"
	"html/template"
	"Coolpy/Account"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	t := template.New("index")
	t, _ = t.ParseFiles("temp/index.html")
	user := &Account.Person{}
	t.Execute(w, user)
}
