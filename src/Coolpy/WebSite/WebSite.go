package WebSite

import (
	"net/http"
	"html/template"
	"fmt"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	t := template.New("index")
	t, _ = t.ParseFiles("temp/index.html")
	t.Execute(w, nil)
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	userName := getUserName(r)
	if userName != "" {
		tmpls, err := compileTemplates("temp/home.html", "temp/homeheader.html", "temp/homefooter.html")
		if err != nil {
			fmt.Fprintf(w, `{"ok":%d,"err":"%v"}`, 0, err)
			return
		}
		home := &Home{Header{Uname:getUserName(r)}}
		tmpls.ExecuteTemplate(w, "home", home)
		tmpls.Execute(w, nil)
	} else {
		http.Redirect(w, r, "/", 302)
	}
}
