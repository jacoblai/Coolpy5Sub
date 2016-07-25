package WebSite

import (
	"net/http"
	"Coolpy/Account"
)

func LoginHandler(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	pass := request.FormValue("password")
	redirectTarget := "/"
	user, err := Account.Get(name)
	if err == nil {
            if user.Pwd == pass{
		    setSession(name, response)
		    redirectTarget = "/home"
	    }
	}
	http.Redirect(response, request, redirectTarget, 302)
}

func LogoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}