package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"Coolpy/Cors"
	"Coolpy/Account"
	"Coolpy/BasicAuth"
	"net"
	"fmt"
	"os"
	"os/signal"
	"Coolpy/Redico"
	"flag"
	"strconv"
)

func main() {
	var port = flag.Int("p", 8080, "coolpy port")

	//初始化数据库服务
	redServer, err := Redico.Run()
	if err != nil {
		panic(err)
	}
	defer redServer.Close()
	redServer.RequireAuth("icoolpy.com")
	//初始化用户账号服务
	Account.Connect(redServer.Addr())
	//自动检测创建超级账号
	Account.CreateAdmin()

	router := httprouter.New()
	router.GET("/:uid", Basicauth.Auth(Account.UserGet))
	router.POST("/", Account.UserPost)

	ln, err := net.Listen("tcp", ":" + strconv.Itoa(*port))
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
			cleanupDone <- true
		}
	}()
	<-cleanupDone
	fmt.Println("\nStoped Coolpy5...\n")
}