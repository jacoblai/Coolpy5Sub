package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"log"
	"Coolpy/Cors"
	"fmt"
	"github.com/satori/go.uuid"
	"Coolpy/Ldata"
	"encoding/json"
)

type Person struct {
	Ukey string
	Uid  string
	Pwd  string
}

func main() {
	ldb := &Ldata.LateEngine{DbPath:"data/ac", DbName:"AccountDb"}
	ldb.Open()
	defer ldb.Ldb.Close()

	np := &Person{
		Ukey : uuid.NewV4().String(),
	}
	np.Uid = "jao"
	np.Pwd = "pwd"

	js, err := json.Marshal(&np)

	err = ldb.Set(np.Uid, []byte(js))
	fmt.Println(err)

	data, err := ldb.Get(np.Uid)
	fmt.Println(string(data), err)

	router := httprouter.New()
	//router.GET("/", Index)
	//router.POST("/", IndexPost)

	if err := http.ListenAndServe(":8080", Cors.CORS(router)); err != nil {
		log.Fatal(err)
	}
}