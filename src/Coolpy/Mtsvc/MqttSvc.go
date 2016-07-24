package Mtsvc

import (
	"github.com/gomqtt/client"
	"github.com/gomqtt/transport"
	"github.com/gomqtt/broker"
	"strconv"
)

var engine *broker.Engine
var mqport int

func Host(mport int) {
	server, err := transport.Launch("tcp://:" + strconv.Itoa(mport))
	if err != nil {
		panic(err)
	}
	engine := broker.NewEngine()
	engine.Accept(server)
	mqport = mport
}

func Close() {
	engine.Close()
}

func Public(k string, payload []byte) {
	client := client.New()
	defer client.Close()
	cf, err := client.Connect("tcp://127.0.0.1:"+ strconv.Itoa(mqport), nil)
	if err != nil {
		panic(err)
	}
	cf.Wait()
	pf, err := client.Publish(k, payload, 0, false)
	if err != nil {
		panic(err)
	}
	pf.Wait()
	err = client.Disconnect()
	if err != nil {
		panic(err)
	}
}