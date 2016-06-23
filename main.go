package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"log"
	"Coolpy/Cors"
	"fmt"
	"Coolpy/Account"
	"encoding/json"
)

type Person struct {
	Ukey string
	Uid  string
	Pwd  string
}

func main() {
	np := Account.New()
	np.Uid = "jao22222"
	np.Pwd = "pwd222"
	np.CreateUkey()
	err := np.CreateOrReplace(np)
	fmt.Println(err)

	data, err := np.Get(np.Uid)
	fmt.Println(data.Ukey, err)

	nnp, err := np.FindKeyStart("jao")
	for k,v := range nnp{
		fmt.Println(k,string(v))
	}



	router := httprouter.New()
	//router.GET("/", Index)
	//router.POST("/", IndexPost)

	if err := http.ListenAndServe(":8080", Cors.CORS(router)); err != nil {
		log.Fatal(err)
	}
}