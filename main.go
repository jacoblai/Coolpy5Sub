package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/surgemq/surgemq/service"
	"github.com/surgemq/surgemq/auth"
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
	"Coolpy/Nodes"
	"Coolpy/Controller"
	"Coolpy/DataPoints"
	"Coolpy/Values"
	"Coolpy/Gpss"
	"Coolpy/Gens"
	"log"
	"Coolpy/MAuth"
)

func main() {
	var (
		port int
		mport int
	)
	flag.IntVar(&port, "p", 8080, "tcp/ip port munber")
	flag.IntVar(&mport, "mp", 1883, "mqtt port munber")
	flag.Parse()
	//初始化数据库服务
	redServer, err := Redico.Run()
	if err != nil {
		panic(err)
	}
	defer redServer.Close()
	svcpwd := "icoolpy.com"
	redServer.RequireAuth(svcpwd)
	//初始化用户账号服务1
	Account.Connect(redServer.Addr(), svcpwd)
	//自动检测创建超级账号
	Account.CreateAdmin()
	//自动id库0
	Incr.Connect(redServer.Addr(), svcpwd)
	//hub库2
	Hubs.Connect(redServer.Addr(), svcpwd)
	//node库3
	Nodes.Connect(redServer.Addr(), svcpwd)
	//控制器库4
	Controller.Connect(redServer.Addr(), svcpwd)
	//数据结点value库5
	Values.Connect(redServer.Addr(), svcpwd)
	//数据结点gps库6
	Gpss.Connect(redServer.Addr(), svcpwd)
	//数据结点gen库7
	Gens.Connect(redServer.Addr(), svcpwd)

	// Create a mqtt server
	auth.Register("coolpy", &MAuth.Manager{})
	mqttsvr := &service.Server{
		KeepAlive:        300, // seconds
		ConnectTimeout:   2, // seconds
		SessionsProvider: "mem", // keeps sessions in memory
		Authenticator:    "coolpy", // always succeed
		TopicsProvider:   "mem", // keeps topic subscriptions in memory
	}
	go func() {
		// Listen and serve connections at mport
		if err := mqttsvr.ListenAndServe("tcp://:" + strconv.Itoa(mport)); err != nil {
			log.Fatal(err)
		}
	}()
	fmt.Println("Coolpy Server mqtt on port ", strconv.Itoa(mport))
	//MqttClient.Connect(strconv.Itoa(mport))

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
	//nodes管理api
	router.POST("/api/hub/:hid/nodes", Basicauth.Auth(Nodes.NodePost))
	router.GET("/api/hub/:hid/nodes", Basicauth.Auth(Nodes.NodesGet))
	router.GET("/api/hub/:hid/node/:nid", Basicauth.Auth(Nodes.NodeGet))
	router.PUT("/api/hub/:hid/node/:nid", Basicauth.Auth(Nodes.NodePut))
	router.DELETE("/api/hub/:hid/node/:nid", Basicauth.Auth(Nodes.NodeDel))
	//datapoints管理api
	router.POST("/api/hub/:hid/node/:nid/datapoints", DataPoints.DPPost)
	router.GET("/api/hub/:hid/node/:nid/datapoint", DataPoints.DPGet)
	router.PUT("/api/hub/:hid/node/:nid/datapoint", DataPoints.DPPut)

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
			mqttsvr.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
	fmt.Println("\nStoped Coolpy5...\n")
}