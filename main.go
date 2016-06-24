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

func main() {
	np := Account.New()
	np.Uid = "li111"
	np.Pwd = "pwd111afasdfasdfasdfasdfasdfasdfqfiweuriquweoruqowieruajsdfkajsdlkfjalskdjfklasjdfklajsdklfjasdklf1122"
	np.CreateUkey()
	err := np.CreateOrReplace(np)
	fmt.Println(err)

	//data, err := np.Get(np.Uid)
	//fmt.Println(data, err)

	//nnp, _ := Account.FindKeyStart("o")
	nnp, _ := Account.All()
	for k,v := range nnp{
		p:= &Account.Person{}
		json.Unmarshal(v,&p)
		fmt.Print(k,p)
	}

	router := httprouter.New()
	//router.GET("/", Index)
	//router.POST("/", IndexPost)

	if err := http.ListenAndServe(":8080", Cors.CORS(router)); err != nil {
		log.Fatal(err)
	}
}