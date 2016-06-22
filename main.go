package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"log"
	"Coolpy/Cors"
	"fmt"
	"Coolpy/Account"
)

type Person struct {
	Ukey string
	Uid  string
	Pwd  string
}

func main() {
	np := Account.New()
	np.Uid = "jao"
	np.Pwd = "pwd111"
	err := np.CreateOrReplace(np)
	fmt.Println(err)

	data, err := np.Get(np.Uid)
	fmt.Println(data.Ukey, err)

	nnp, err := np.Find(np.Uid)
	fmt.Println(string(nnp[np.Uid]))


	router := httprouter.New()
	//router.GET("/", Index)
	//router.POST("/", IndexPost)

	if err := http.ListenAndServe(":8080", Cors.CORS(router)); err != nil {
		log.Fatal(err)
	}
}