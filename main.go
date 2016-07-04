package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"Coolpy/Cors"
	"Coolpy/Account"
	"encoding/json"
	"Coolpy/BasicAuth"
	"net"
	"fmt"
	"Coolpy/Incr"
	"os"
	"os/signal"
)

func main() {
	router := httprouter.New()
	router.GET("/:uid", Basicauth.Auth(Index))
	router.POST("/", IndexPost)

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Can't listen: %s", err)
	}
	go http.Serve(ln, Cors.CORS(router))
	fmt.Println("Coolpy Server host on port 8080")
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			fmt.Println("\nStopping Coolpy5...\n")
			ln.Close()
			Account.Close()
			Incr.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
	fmt.Println("\nStoped Coolpy5...\n")
}

func Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	p, err := Account.Get(ps.ByName("uid"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(p)
	return
}

func IndexPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var p Account.Person
	err := decoder.Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println(err)
		return
	}
	p.CreateUkey()
	err = Account.CreateOrReplace(&p)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println(err)
		return
	}
	json.NewEncoder(w).Encode(p)
	return
}