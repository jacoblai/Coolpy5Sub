package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"log"
	"Coolpy/Cors"
	"Coolpy/Account"
	"fmt"
)

func main() {
	np := Account.New()
	np.Uid = "jao"
	np.Pwd = "pwd"

	err := np.Save(np)
	fmt.Println(err)

	router := httprouter.New()
	//router.GET("/", Index)
	//router.POST("/", IndexPost)

	if err := http.ListenAndServe(":8080", Cors.CORS(router)); err != nil {
		log.Fatal(err)
	}
}