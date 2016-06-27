package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"log"
	"Coolpy/Cors"
	"Coolpy/Account"
	"encoding/json"
	"Coolpy/BasicAuth"
	"os"
	"os/signal"
	"net"
)

func main() {
	router := httprouter.New()
	router.GET("/:uid", Basicauth.Auth(Index))
	router.POST("/", IndexPost)

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Can't listen: %s", err)
	}
	go http.Serve(ln, Cors.CORS(router))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	log.Println("Coolpy server on stopped signal is:", s)
	ln.Close()
	os.Exit(1)
}

func Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	p, err := Account.Get(ps.ByName("uid"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(p)
}

func IndexPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var p Account.Person
	err := decoder.Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Fatal(err)
		return
	}
	p.CreateUkey()
	err = Account.CreateOrReplace(&p)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Fatal(err)
		return
	}
	json.NewEncoder(w).Encode(p)
}