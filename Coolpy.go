package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"Coolpy/Cors"
	"net"
	"fmt"
	"os"
	"os/signal"
	"Coolpy/Redico"
	"flag"
	"strconv"
	"Coolpy/Mtsvc"
	"log"
	"Coolpy/CoSystem"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"Coolpy"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("Coolpy Version:", CoSystem.CpVersion)
	var (
		port = flag.Int("a", 6543, "web api port munber")
		mport = flag.Int("m", 1883, "mqtt port munber")
		wsport = flag.Int("s", 1884, "mqtt websocket port munber")
		wport = flag.Int("w", 8000, "www website port munber")
	)
	flag.Parse()

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	//初始化数据库服务
	redServer, err := Redico.Run(dir)
	if err != nil {
		panic(err)
	}
	defer redServer.Close()
	svcpwd := "icoolpy.com"
	redServer.RequireAuth(svcpwd)
	//初始化用户账号服务1
	Coolpy.AccConnect(redServer.Addr(), svcpwd)
	//自动检测创建超级账号
	Coolpy.CreateAdmin()
	//自动id库0
	Coolpy.InrcConnect(redServer.Addr(), svcpwd)
	//hub库2
	Coolpy.HubConnect(redServer.Addr(), svcpwd)
	//node库3
	Coolpy.NodeConnect(redServer.Addr(), svcpwd)
	//控制器库4
	Coolpy.CtrlConnect(redServer.Addr(), svcpwd)
	//数据结点value库5
	Coolpy.ValdpConnect(redServer.Addr(), svcpwd)
	//数据结点gps库6
	Coolpy.GpsdpConnect(redServer.Addr(), svcpwd)
	//数据结点gen库7
	Coolpy.GendpConnect(redServer.Addr(), svcpwd)
	//数据结点img库8
	Coolpy.PhotoConnect(redServer.Addr(), svcpwd)

	//host mqtt service
	go func() {
		msvc := &Mtsvc.MqttSvc{}
		msvc.Host(*mport, *wsport)
	}()
	fmt.Println("Coolpy mqtt on port", strconv.Itoa(*mport))
	fmt.Println("Coolpy mqtt websocket on port", strconv.Itoa(*wsport))
	router := httprouter.New()
	//用户管理api
	router.POST("/api/user", Coolpy.Auth(Coolpy.UserPost))
	router.GET("/api/user/:uid", Coolpy.Auth(Coolpy.UserGet))
	router.PUT("/api/user/:uid", Coolpy.Auth(Coolpy.UserPut))
	router.DELETE("/api/user/:uid", Coolpy.Auth(Coolpy.UserDel))
	router.GET("/api/um/all", Coolpy.Auth(Coolpy.UserAll))
	router.GET("/api/um/apikey", Coolpy.Auth(Coolpy.UserApiKey))
	router.GET("/api/um/newkey", Coolpy.Auth(Coolpy.UserNewApiKey))
	//hubs管理api
	router.POST("/api/hubs", Coolpy.Auth(Coolpy.HubPost))
	router.GET("/api/hubs", Coolpy.Auth(Coolpy.HubsGet))
	router.GET("/api/hubs/all", Coolpy.Auth(Coolpy.HubsAll))
	router.GET("/api/hub/:hid", Coolpy.Auth(Coolpy.HubGet))
	router.PUT("/api/hub/:hid", Coolpy.Auth(Coolpy.HubPut))
	router.DELETE("/api/hub/:hid", Coolpy.Auth(Coolpy.HubDel))
	//nodes管理api
	router.POST("/api/hub/:hid/nodes", Coolpy.Auth(Coolpy.NodePost))
	router.GET("/api/hub/:hid/nodes", Coolpy.Auth(Coolpy.NodesGet))
	router.GET("/api/hub/:hid/node/:nid", Coolpy.Auth(Coolpy.NodeGet))
	router.PUT("/api/hub/:hid/node/:nid", Coolpy.Auth(Coolpy.NodePut))
	router.DELETE("/api/hub/:hid/node/:nid", Coolpy.Auth(Coolpy.NodeDel))
	//datapoints管理api
	router.POST("/api/hub/:hid/node/:nid/datapoints", Coolpy.DPPost)//传感器提交单个数据结点
	router.GET("/api/hub/:hid/node/:nid/datapoint", Coolpy.DPGet)//所有控制器及传感器取最新值
	router.PUT("/api/hub/:hid/node/:nid/datapoint", Coolpy.DPPut)//控制器更新值
	router.GET("/api/hub/:hid/node/:nid/datapoint/:key", Coolpy.DPGetByKey)//传感器取得key对应值
	router.PUT("/api/hub/:hid/node/:nid/datapoint/:key", Coolpy.DPPutByKey)//传感器更新key对应值
	router.DELETE("/api/hub/:hid/node/:nid/datapoint/:key", Coolpy.DPDelByKey)//传感器删除key对应值
	router.GET("/api/hub/:hid/node/:nid/json", Coolpy.DPGetRange)//传感器取得历史数据
	//图像管理api
	router.POST("/api/hub/:hid/node/:nid/photos", Coolpy.PhotoPost)//上传图片png,jpg,gif
	router.GET("/api/hub/:hid/node/:nid/photo/content", Coolpy.PhotoGet)
	router.GET("/api/hub/:hid/node/:nid/photo/content/:key", Coolpy.PhotoGetByKey)
	router.DELETE("/api/hub/:hid/node/:nid/photo/content/:key", Coolpy.PhotoDelByKey)
	//系统api
	router.GET("/api/sys/version", CoSystem.VersionGet)
	//系统底层api
	router.POST("/os/cmd", Coolpy.Auth(Coolpy.CmdPost))
	router.POST("/os/upload/:filename", Coolpy.Auth(Coolpy.UploadPost))

	go func() {
		ln, err := net.Listen("tcp", ":" + strconv.Itoa(*port))
		if err != nil {
			fmt.Println("Can't listen:", err)
		}
		err = http.Serve(ln, Cors.CORS(router))
		if err != nil {
			log.Fatal(err)
		}
	}()
	fmt.Println("Coolpy http on port", strconv.Itoa(*port))

	//当api端口号被启动参数修改时即自动更新www相关连接参数
	settingpath := dir + "/www/scripts-app/setting.js"
	nstring := "var basicurl=\"http://\"+ window.location.hostname +\":" + strconv.Itoa(*port) + "\""
	err = ioutil.WriteFile(settingpath, []byte(nstring), 0644)
	if err != nil {
		fmt.Println(err)
	}
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(dir + "/www"))))
	go func() {
		err := http.ListenAndServe(":" + strconv.Itoa(*wport), nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
	fmt.Println("Coolpy www on port", strconv.Itoa(*wport))
	fmt.Println("Power By ICOOLPY.COM")

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			cleanupDone <- true
		}
	}()
	<-cleanupDone
	fmt.Println("\nStoped Coolpy5...\n")
}