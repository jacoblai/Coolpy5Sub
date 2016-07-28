package Mtsvc

import (
	"github.com/gomqtt/client"
	"github.com/gomqtt/transport"
	"github.com/gomqtt/broker"
	"strconv"
	"fmt"
)

type MqttSvc struct {
	Engine *broker.Engine
}

var Mport int

func (m *MqttSvc) Host(mport int) {
	Mport = mport
	server, err := transport.Launch("tcp://:" + strconv.Itoa(Mport))
	if err != nil {
		panic(err)
	}
	engine := broker.NewEngine()
	engine.Accept(server)
	m.Engine = broker.NewEngine()
	m.Engine.Accept(server)
}

func Public(k string, payload []byte) {
	client := client.New()
	defer client.Close()
	cf, err := client.Connect("tcp://127.0.0.1:"+ strconv.Itoa(Mport), nil)
	if err != nil {
		fmt.Println(err)
	}
	cf.Wait()
	pf, err := client.Publish(k, payload, 0, false)
	if err != nil {
		fmt.Println(err)
	}
	pf.Wait()
	err = client.Disconnect()
	if err != nil {
		fmt.Println(err)
	}
}