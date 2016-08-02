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
	"Coolpy/Nodes"
	"Coolpy/Controller"
	"Coolpy/DataPoints"
	"Coolpy/Values"
	"Coolpy/Gpss"
	"Coolpy/Gens"
	"Coolpy/Mtsvc"
	"Coolpy/Photos"
	"log"
	"Coolpy/Deller"
	"Coolpy/CoSystem"
	"io/ioutil"
	"strings"
)

func main() {
	fmt.Println("Coolpy Version:", CoSystem.CpVersion)
	var (
		port int
		mport int
		wport int
	)
	flag.IntVar(&port, "a", 6543, "web api port munber")
	flag.IntVar(&mport, "m", 1883, "mqtt port munber")
	flag.IntVar(&wport, "w", 8000, "www website port munber")
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
	//数据结点img库8
	Photos.Connect(redServer.Addr(), svcpwd)

	//host mqtt service
	go func() {
		msvc := &Mtsvc.MqttSvc{}
		msvc.Host(mport)
		defer msvc.Engine.Close()
	}()
	fmt.Println("Coolpy mqtt on port", strconv.Itoa(mport))

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
	router.GET("/api/hubs", Basicauth.Auth(Hubs.HubsGet))
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
	router.POST("/api/hub/:hid/node/:nid/datapoints", DataPoints.DPPost)//传感器提交单个数据结点
	router.GET("/api/hub/:hid/node/:nid/datapoint", DataPoints.DPGet)//所有控制器及传感器取最新值
	router.PUT("/api/hub/:hid/node/:nid/datapoint", DataPoints.DPPut)//控制器更新值
	router.GET("/api/hub/:hid/node/:nid/datapoint/:key", DataPoints.DPGetByKey)//传感器取得key对应值
	router.PUT("/api/hub/:hid/node/:nid/datapoint/:key", DataPoints.DPPutByKey)//传感器更新key对应值
	router.DELETE("/api/hub/:hid/node/:nid/datapoint/:key", DataPoints.DPDelByKey)//传感器删除key对应值
	router.GET("/api/hub/:hid/node/:nid/json", DataPoints.DPGetRange)//传感器取得历史数据
	//图像管理api
	router.POST("/api/hub/:hid/node/:nid/photos", Photos.PhotoPost)//上传图片png,jpg,gif
	router.GET("/api/hub/:hid/node/:nid/photo/content", Photos.PhotoGet)
	router.GET("/api/hub/:hid/node/:nid/photo/content/:key", Photos.PhotoGetByKey)
	router.DELETE("/api/hub/:hid/node/:nid/photo/content/:key", Photos.PhotoDelByKey)
	//系统api
	router.GET("/api/sys/version", CoSystem.VersionGet)
	go func() {
		ln, err := net.Listen("tcp", ":" + strconv.Itoa(port))
		if err != nil {
			fmt.Println("Can't listen:", err)
		}
		err = http.Serve(ln, Cors.CORS(router))
		if err != nil {
			log.Fatal(err)
		}
	}()
	fmt.Println("Coolpy http on port", strconv.Itoa(port))
	if port != 6543 {
		//当api端口号被启动参数修改时即自动更新www相关连接参数
		f, err := ioutil.ReadFile("www/scripts-app/setting.js")
		if err != nil {
			fmt.Println(err)
		}
		nstring := strings.Replace(string(f), "6543", strconv.Itoa(port), -1)
		err = ioutil.WriteFile("www/scripts-app/setting.js", []byte(nstring), 0644)
		if err != nil{
			fmt.Println(err)
		}
	}

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("www"))))
	go func() {
		err := http.ListenAndServe(":" + strconv.Itoa(wport), nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
	fmt.Println("Coolpy www on port", strconv.Itoa(wport))
	fmt.Println("Power By ICOOLPY.COM")

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			close(Deller.DelControl)
			close(Deller.DelControls)
			close(Deller.DelHub)
			close(Deller.DelHubs)
			close(Deller.DelNode)
			close(Deller.DelNodes)
			close(Deller.DelValues)
			close(Deller.DelGpss)
			close(Deller.DelGens)
			close(Deller.DelPhotos)
			cleanupDone <- true
		}
	}()
	<-cleanupDone
	fmt.Println("\nStoped Coolpy5...\n")
}