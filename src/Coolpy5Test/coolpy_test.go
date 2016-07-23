package Coolpy5Test

import (
	"testing"
	"Coolpy/Hubs"
	"Coolpy/Account"
	"Coolpy/Nodes"
	"Coolpy/Controller"
	"Coolpy/Values"
	"Coolpy/Gpss"
	"Coolpy/Gens"
	"Coolpy/Photos"
	"fmt"
)

func TestCoolpy5(t *testing.T) {
	Addr := "127.0.0.1:6380"
	svcpwd := "icoolpy.com"
	//初始化用户账号服务1
	Account.Connect(Addr, svcpwd)
	//hub库2
	Hubs.Connect(Addr, svcpwd)
	//node库3
	Nodes.Connect(Addr, svcpwd)
	//控制器库4
	Controller.Connect(Addr, svcpwd)
	//数据结点value库5
	Values.Connect(Addr, svcpwd)
	//数据结点gps库6
	Gpss.Connect(Addr, svcpwd)
	//数据结点gen库7
	Gens.Connect(Addr, svcpwd)
	//数据结点img库8
	Photos.Connect(Addr, svcpwd)

	acs,err := Account.All()
	if err !=nil{
		t.Error(err)
	}
	for _,v:=range acs{
		fmt.Println("users",v)
	}

	hubs,err := Hubs.All()
	if err !=nil{
		t.Error(err)
	}
	for _,v:= range hubs{
		fmt.Println("hubs",v)
	}

	nodes,err := Nodes.All()
	if err !=nil{
		t.Error(err)
	}
	for _,v:= range nodes{
		fmt.Println("nodes",v)
	}

	Control,err := Controller.All()
	if err !=nil{
		t.Error(err)
	}
	for _,v:= range Control{
		fmt.Println("Control",v)
	}

	vals,err := Values.All()
	if err !=nil{
		t.Error(err)
	}
	for _,v:= range vals{
		fmt.Println("vals",v)
	}

	gpss,err := Gpss.All()
	if err !=nil{
		t.Error(err)
	}
	for _,v:= range gpss{
		fmt.Println("gpss",v)
	}

	gens,err := Gens.All()
	if err !=nil{
		t.Error(err)
	}
	for _,v:= range gens{
		fmt.Println("gens",v)
	}

	photos,err := Photos.All()
	if err !=nil{
		t.Error(err)
	}
	for _,v:= range photos{
		fmt.Println("photos",v)
	}

	//req, _ := http.NewRequest("GET", "http://localhost:8080/api/hub/5/node/18/photo/content", nil)
	//req.Header.Add("U-ApiKey","83f68172-62a7-412d-96b6-5db3f15c3983")
	//req.Header.Add("Range", "bytes=1023-")
	//var client http.Client
	//resp, _ := client.Do(req)
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(len(body))
}