package Coolpy5Test

import (
	"testing"
	"fmt"
	"Coolpy"
)

func TestCoolpy5(t *testing.T) {
	Addr := "127.0.0.1:6380"
	svcpwd := "icoolpy.com"
	//初始化用户账号服务1
	Coolpy.AccConnect(Addr, svcpwd)
	Coolpy.HubConnect(Addr, svcpwd)
	//node库3
	Coolpy.NodeConnect(Addr, svcpwd)
	//控制器库4
	Coolpy.CtrlConnect(Addr, svcpwd)
	//数据结点value库5
	Coolpy.ValdpConnect(Addr, svcpwd)
	//数据结点gps库6
	Coolpy.GpsdpConnect(Addr, svcpwd)
	//数据结点gen库7
	Coolpy.GendpConnect(Addr, svcpwd)
	//数据结点img库8
	Coolpy.PhotoConnect(Addr, svcpwd)

	acs,err := Coolpy.AccAll()
	if err !=nil{
		t.Error(err)
	}
	for _,v:=range acs{
		fmt.Println("users",v)
	}

	hubs,err := Coolpy.HubAll()
	if err !=nil{
		t.Error(err)
	}
	for _,v:= range hubs{
		fmt.Println("hubs",v)
	}

	nodes,err := Coolpy.NodeAll()
	if err !=nil{
		t.Error(err)
	}
	for _,v:= range nodes{
		fmt.Println("nodes",v)
	}

	Control,err := Coolpy.CtrlAll()
	if err !=nil{
		t.Error(err)
	}
	for _,v:= range Control{
		fmt.Println("Control",v)
	}

	vals,err := Coolpy.ValdpAll()
	if err !=nil{
		t.Error(err)
	}
	for _,v:= range vals{
		fmt.Println("vals",v)
	}

	gpss,err := Coolpy.GpsdpAll()
	if err !=nil{
		t.Error(err)
	}
	for _,v:= range gpss{
		fmt.Println("gpss",v)
	}

	gens,err := Coolpy.GendpAll()
	if err !=nil{
		t.Error(err)
	}
	for _,v:= range gens{
		fmt.Println("gens",v)
	}

	photos,err := Coolpy.PhotoAll()
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