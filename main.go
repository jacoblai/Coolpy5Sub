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
	"Coolpy/Incr"
	"Coolpy/Hubs"
)

func main() {
	var port int
	flag.IntVar(&port, "p", 8080, "tcp/ip port munber")
	flag.Parse()
	//初始化数据库服务
	redServer, err := Redico.Run()
	if err != nil {
		panic(err)
	}
	defer redServer.Close()
	svcpwd := "icoolpy.com"
	redServer.RequireAuth(svcpwd)
	//初始化用户账号服务
	Account.Connect(redServer.Addr(), svcpwd)
	//自动检测创建超级账号
	Account.CreateAdmin()
	//自动id库
	Incr.Connect(redServer.Addr(), svcpwd)
	//hub库
	Hubs.Connect(redServer.Addr(), svcpwd)

	router := httprouter.New()
	//用户管理api
	router.POST("/api/user", Basicauth.Auth(Account.UserPost))
	router.GET("/api/user/:uid", Basicauth.Auth(Account.UserGet))
	router.PUT("/api/user/:uid", Basicauth.Auth(Account.UserPut))
	router.DELETE("/api/user/:uid", Basicauth.Auth(Account.UserDel))
	router.GET("/api/um/all", Basicauth.Auth(Account.UserAll))
	router.GET("/api/um/apikey", Basicauth.Auth(Account.UserApiKey))
	router.GET("/api/um/newkey", Basicauth.Auth(Account.UserNewApiKey))
	//hubs管理api
	router.POST("/api/hubs", Basicauth.Auth(Hubs.HubPost))
	router.GET("/api/hubs/:ukey", Basicauth.Auth(Hubs.HubsGet))
	router.GET("/api/hub/:hid", Basicauth.Auth(Hubs.HubGet))
	router.PUT("/api/hub/:hid", Basicauth.Auth(Hubs.HubPut))
	router.DELETE("/api/hub/:hid", Basicauth.Auth(Hubs.HubDel))

	ln, err := net.Listen("tcp", ":" + strconv.Itoa(port))
	if err != nil {
		fmt.Println("Can't listen: %s", err)
	}
	go http.Serve(ln, Cors.CORS(router))
	fmt.Println("Coolpy Server host on port ", strconv.Itoa(port))

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